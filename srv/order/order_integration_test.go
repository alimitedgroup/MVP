package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/order/adapter"
	"github.com/alimitedgroup/MVP/srv/order/adapter/controller"
	"github.com/alimitedgroup/MVP/srv/order/adapter/listener"
	"github.com/alimitedgroup/MVP/srv/order/business"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/zap/zaptest"
)

var Module = fx.Options(
	lib.Module,
	adapter.Module,
	business.Module,
)

type TestParams struct {
	Js jetstream.JetStream
	Ns *nats.Conn
}

func IntegrationTest(t *testing.T, testFunc any) {
	ctx := t.Context()

	ns, _ := broker.NewInProcessNATSServer(t)

	js, err := jetstream.New(ns)
	require.NoError(t, err)

	p := TestParams{
		Ns: ns,
		Js: js,
	}

	app := fx.New(
		Module,
		fx.Supply(&p),
		fx.Supply(ns),
		fx.Invoke(testFunc),
		fx.Supply(zaptest.NewLogger(t)),
	)
	err = app.Start(ctx)
	require.NoError(t, err)

	defer func() {
		err := app.Stop(ctx)
		require.NoError(t, err)
	}()
}

func NatsSendRequest[RequestDTO any, ResponseDTO any](t *testing.T, p *TestParams, subject string, dto RequestDTO) ResponseDTO {
	payload, err := json.Marshal(dto)
	require.NoError(t, err)

	msg, err := p.Ns.Request(subject, payload, 1*time.Second)
	require.NoError(t, err)

	var respDto ResponseDTO
	err = json.Unmarshal(msg.Data, &respDto)
	require.NoError(t, err)

	return respDto
}

func NatsPublishToStream[RequestDTO any](ctx context.Context, t *testing.T, p *TestParams, stream jetstream.StreamConfig, subject string, event RequestDTO) {
	s, err := p.Js.CreateStream(ctx, stream)
	require.NoError(t, err)

	beforeInfo, err := s.Info(ctx)
	require.NoError(t, err)

	payload, err := json.Marshal(event)
	require.NoError(t, err)

	ack, err := p.Js.Publish(ctx, subject, payload)
	require.NoError(t, err)

	info, err := s.Info(ctx)
	require.NoError(t, err)

	require.Equal(t, info.State.LastSeq, beforeInfo.State.LastSeq+1)
	require.Equal(t, ack.Stream, stream.Name)
}

func TestCreateOrder(t *testing.T) {
	IntegrationTest(t, func(
		lc fx.Lifecycle, p *TestParams,
		orderControllerRouter *controller.OrderRouter,
		transferControllerRouter *controller.TransferRouter,
		orderListenerRouter *listener.OrderRouter,
		stockUpdateListenerRouter *listener.StockUpdateRouter,
	) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				for _, v := range []lib.BrokerRoute{orderControllerRouter, transferControllerRouter, orderListenerRouter, stockUpdateListenerRouter} {
					err := v.Setup(ctx)
					require.NoError(t, err)
				}

				sub, err := p.Ns.Subscribe("warehouse.1.reservation.create", func(msg *nats.Msg) {
					var event request.ReserveStockRequestDTO
					err := json.Unmarshal(msg.Data, &event)
					require.NoError(t, err)

					require.Equal(t, event.Goods[0].GoodID, "1")
					require.Equal(t, event.Goods[0].Quantity, int64(5))

					err = msg.Respond([]byte(`{"error": "", "message": {"reservation_id": "1"}}`))
					require.NoError(t, err)
				})
				require.NoError(t, err)
				defer func() {
					err := sub.Unsubscribe()
					require.NoError(t, err)
				}()

				{
					event := stream.StockUpdate{
						ID:          "1",
						WarehouseID: "1",
						Type:        stream.StockUpdateTypeAdd,
						Goods: []stream.StockUpdateGood{
							{GoodID: "1", Quantity: 10, Delta: 10},
						},
						OrderID:       "",
						TransferID:    "",
						ReservationID: "",
						Timestamp:     time.Now().UnixMilli(),
					}

					NatsPublishToStream(ctx, t, p, stream.StockUpdateStreamConfig, "stock.update.1", event)
				}

				var orderId string
				{
					dto := request.CreateOrderRequestDTO{
						Name:     "John Doe Order",
						FullName: "John Doe",
						Address:  "123 Main St",
						Goods: []request.CreateOrderGood{
							{GoodID: "1", Quantity: 5},
						},
					}
					resp := NatsSendRequest[request.CreateOrderRequestDTO, response.OrderCreateResponseDTO](
						t, p, "order.create", dto,
					)

					require.Empty(t, resp.Error)
					require.NotEmpty(t, resp.Message.OrderID)
					orderId = resp.Message.OrderID
				}

				// don't wait to get the Created state
				time.Sleep(0 * time.Millisecond)

				{
					resp := NatsSendRequest[any, response.GetAllOrderResponseDTO](
						t, p, "order.get.all", map[string]interface{}{},
					)

					require.Empty(t, resp.Error)
					require.Equal(t, orderId, resp.Message[0].OrderID)
					require.Equal(t, "Created", resp.Message[0].Status)
				}

				// wait for background processing of the order
				time.Sleep(10 * time.Millisecond)

				{
					resp := NatsSendRequest[any, response.GetAllOrderResponseDTO](
						t, p, "order.get.all", map[string]interface{}{},
					)

					require.Empty(t, resp.Error)
					require.Equal(t, orderId, resp.Message[0].OrderID)
					require.Equal(t, "1", resp.Message[0].Reservations[0])
					require.Equal(t, "Filled", resp.Message[0].Status)
				}

				time.Sleep(10 * time.Millisecond)

				{
					event := stream.StockUpdate{
						ID:          "2",
						WarehouseID: "1",
						Type:        stream.StockUpdateTypeOrder,
						Goods: []stream.StockUpdateGood{
							{GoodID: "1", Quantity: 5, Delta: 5},
						},
						OrderID:       orderId,
						TransferID:    "",
						ReservationID: "1",
						Timestamp:     time.Now().UnixMilli(),
					}

					NatsPublishToStream(ctx, t, p, stream.StockUpdateStreamConfig, "stock.update.1", event)
				}

				time.Sleep(10 * time.Millisecond)

				{
					resp := NatsSendRequest[any, response.GetAllOrderResponseDTO](
						t, p, "order.get.all", map[string]interface{}{},
					)

					require.Empty(t, resp.Error)
					require.Equal(t, orderId, resp.Message[0].OrderID)
					require.Equal(t, "1", resp.Message[0].Reservations[0])
					require.Equal(t, "Completed", resp.Message[0].Status)

				}

				return nil
			},
		})

	})
}

func TestCreateTransfer(t *testing.T) {
	IntegrationTest(t, func(
		lc fx.Lifecycle, p *TestParams,
		orderControllerRouter *controller.OrderRouter,
		transferControllerRouter *controller.TransferRouter,
		orderListenerRouter *listener.OrderRouter,
		stockUpdateListenerRouter *listener.StockUpdateRouter,
	) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				for _, v := range []lib.BrokerRoute{orderControllerRouter, transferControllerRouter, orderListenerRouter, stockUpdateListenerRouter} {
					err := v.Setup(ctx)
					require.NoError(t, err)
				}

				sub, err := p.Ns.Subscribe("warehouse.1.reservation.create", func(msg *nats.Msg) {
					var event request.ReserveStockRequestDTO
					err := json.Unmarshal(msg.Data, &event)
					require.NoError(t, err)

					require.Equal(t, event.Goods[0].GoodID, "1")
					require.Equal(t, event.Goods[0].Quantity, int64(5))

					err = msg.Respond([]byte(`{"error": "", "message": {"reservation_id": "1"}}`))
					require.NoError(t, err)
				})
				require.NoError(t, err)
				defer func() {
					err := sub.Unsubscribe()
					require.NoError(t, err)
				}()

				{
					event := stream.StockUpdate{
						ID:          "1_1",
						WarehouseID: "1",
						Type:        stream.StockUpdateTypeAdd,
						Goods: []stream.StockUpdateGood{
							{GoodID: "1", Quantity: 10, Delta: 10},
						},
						OrderID:       "",
						TransferID:    "",
						ReservationID: "",
						Timestamp:     time.Now().UnixMilli(),
					}

					NatsPublishToStream(ctx, t, p, stream.StockUpdateStreamConfig, "stock.update.1", event)
				}

				var transferId string
				{
					dto := request.CreateTransferRequestDTO{
						SenderID:   "1",
						ReceiverID: "2",
						Goods: []request.TransferGood{
							{GoodID: "1", Quantity: 5},
						},
					}
					resp := NatsSendRequest[request.CreateTransferRequestDTO, response.TransferCreateResponseDTO](
						t, p, "transfer.create", dto,
					)

					require.Empty(t, resp.Error)
					require.NotEmpty(t, resp.Message.TransferID)
					transferId = resp.Message.TransferID
				}

				// don't wait to get the Created state
				time.Sleep(0 * time.Millisecond)

				{
					resp := NatsSendRequest[any, response.GetAllTransferResponseDTO](
						t, p, "transfer.get.all", map[string]interface{}{},
					)

					require.Empty(t, resp.Error)
					require.Equal(t, transferId, resp.Message[0].TransferID)
					require.Equal(t, "Created", resp.Message[0].Status)
				}

				// wait for background processing of the order
				time.Sleep(10 * time.Millisecond)

				{
					resp := NatsSendRequest[any, response.GetAllTransferResponseDTO](
						t, p, "transfer.get.all", map[string]interface{}{},
					)

					require.Empty(t, resp.Error)
					require.Equal(t, transferId, resp.Message[0].TransferID)
					require.Equal(t, "Filled", resp.Message[0].Status)
				}

				time.Sleep(10 * time.Millisecond)

				{
					event := stream.StockUpdate{
						ID:          "1_2",
						WarehouseID: "1",
						Type:        stream.StockUpdateTypeTransfer,
						Goods: []stream.StockUpdateGood{
							{GoodID: "1", Quantity: 5, Delta: 5},
						},
						OrderID:       "",
						TransferID:    transferId,
						ReservationID: "1",
						Timestamp:     time.Now().UnixMilli(),
					}

					NatsPublishToStream(ctx, t, p, stream.StockUpdateStreamConfig, "stock.update.1", event)
				}

				{
					event := stream.StockUpdate{
						ID:          "2_1",
						WarehouseID: "2",
						Type:        stream.StockUpdateTypeTransfer,
						Goods: []stream.StockUpdateGood{
							{GoodID: "1", Quantity: 5, Delta: 5},
						},
						OrderID:       "",
						TransferID:    transferId,
						ReservationID: "1",
						Timestamp:     time.Now().UnixMilli(),
					}

					NatsPublishToStream(ctx, t, p, stream.StockUpdateStreamConfig, "stock.update.2", event)
				}

				time.Sleep(10 * time.Millisecond)

				{
					resp := NatsSendRequest[any, response.GetAllTransferResponseDTO](
						t, p, "transfer.get.all", map[string]interface{}{},
					)

					require.Empty(t, resp.Error)
					require.Equal(t, transferId, resp.Message[0].TransferID)
					require.Equal(t, "Completed", resp.Message[0].Status)
				}

				return nil
			},
		})

	})
}
