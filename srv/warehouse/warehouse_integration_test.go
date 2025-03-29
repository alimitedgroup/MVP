package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/controller"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/listener"
	"github.com/alimitedgroup/MVP/srv/warehouse/business"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

var modules = fx.Options(
	lib.ModuleTest,
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
	cfg := config.WarehouseConfig{ID: "1"}

	js, err := jetstream.New(ns)
	require.NoError(t, err)

	p := TestParams{
		Ns: ns,
		Js: js,
	}

	app := fx.New(
		modules,
		fx.Supply(ns, &p, &cfg, t),
		fx.Invoke(testFunc),
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

func TestAddAndRemoveWarehouseStock(t *testing.T) {
	IntegrationTest(t, func(
		lc fx.Lifecycle, cfg *config.WarehouseConfig, p *TestParams,
		stockControllerRouter *controller.StockRouter, catalogListenerRouter *listener.CatalogRouter, stockListenerRouter *listener.StockUpdateRouter) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				for _, v := range []lib.BrokerRoute{catalogListenerRouter, stockListenerRouter, stockControllerRouter} {
					err := v.Setup(ctx)
					require.NoError(t, err)
				}

				{
					event := stream.GoodUpdateData{
						GoodID:             "1",
						GoodNewName:        "hat",
						GoodNewDescription: "very good hat",
					}

					NatsPublishToStream(ctx, t, p, stream.AddOrChangeGoodDataStream, "good.update", event)
					time.Sleep(10 * time.Millisecond)
				}

				{
					dto := request.AddStockRequestDTO{GoodID: "1", Quantity: 10}
					resp := NatsSendRequest[request.AddStockRequestDTO, response.ResponseDTO[string]](
						t, p, fmt.Sprintf("warehouse.%s.stock.add", cfg.ID), dto,
					)

					require.Empty(t, resp.Error)
					require.Equal(t, resp.Message, "ok")
				}

				time.Sleep(10 * time.Millisecond)

				{
					dto := request.RemoveStockRequestDTO{GoodID: "1", Quantity: 5}
					resp := NatsSendRequest[request.RemoveStockRequestDTO, response.ResponseDTO[string]](
						t, p, fmt.Sprintf("warehouse.%s.stock.remove", cfg.ID), dto,
					)

					require.Empty(t, resp.Error)
					require.Equal(t, resp.Message, "ok")
				}

				time.Sleep(10 * time.Millisecond)

				{
					s, err := p.Js.Stream(ctx, stream.StockUpdateStreamConfig.Name)
					require.NoError(t, err)

					msg, err := s.GetLastMsgForSubject(ctx, fmt.Sprintf("stock.update.%s", cfg.ID))
					require.NoError(t, err)

					var event stream.StockUpdate
					err = json.Unmarshal(msg.Data, &event)
					require.NoError(t, err)

					require.Equal(t, event.Goods[0].GoodID, "1")
					require.Equal(t, event.Goods[0].Quantity, int64(5))
					require.Equal(t, event.WarehouseID, cfg.ID)
					require.Equal(t, event.Type, stream.StockUpdateTypeRemove)
				}

				return nil
			},
		})
	})
}
