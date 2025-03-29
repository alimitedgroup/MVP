package controller

import (
	"cmp"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/lib"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
	"golang.org/x/exp/slices"
)

func TestOrderControllerCreateOrder(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)

	createOrderUseCaseMock := NewMockICreateOrderUseCase(ctrl)
	createOrderUseCaseMock.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(port.CreateOrderResponse{OrderID: "1"}, nil)

	getOrderUseCaseMock := NewMockIGetOrderUseCase(ctrl)

	ns, _ := broker.NewInProcessNATSServer(t)

	app := fx.New(
		lib.ModuleTest,
		fx.Supply(ns, t),
		fx.Supply(fx.Annotate(createOrderUseCaseMock, fx.As(new(port.ICreateOrderUseCase)))),
		fx.Supply(fx.Annotate(getOrderUseCaseMock, fx.As(new(port.IGetOrderUseCase)))),
		fx.Provide(NewOrderController),
		fx.Provide(NewOrderRouter),
		fx.Invoke(func(lc fx.Lifecycle, r *OrderRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					require.NoError(t, err)

					dto := request.CreateOrderRequestDTO{
						Name:     "Order 1",
						FullName: "test test",
						Address:  "via roma 11",
						Goods: []request.CreateOrderGood{
							{
								GoodID:   "1",
								Quantity: 10,
							},
						},
					}
					payload, err := json.Marshal(dto)
					require.NoError(t, err)

					resp, err := ns.Request("order.create", payload, 1*time.Second)
					require.NoError(t, err)

					var respDto response.OrderCreateResponseDTO
					err = json.Unmarshal(resp.Data, &respDto)
					require.NoError(t, err)

					require.Empty(t, respDto.Error)
					require.Equal(t, respDto.Message.OrderID, "1")

					return nil
				},
			})
		}),
	)

	err := app.Start(ctx)
	require.NoError(t, err)

	defer func() {
		err = app.Stop(ctx)
		require.NoError(t, err)
	}()
}

func TestOrderControllerGetOrder(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)

	createOrderUseCaseMock := NewMockICreateOrderUseCase(ctrl)

	getOrderUseCaseMock := NewMockIGetOrderUseCase(ctrl)
	getOrderUseCaseMock.EXPECT().GetOrder(gomock.Any(), gomock.Any()).Return(model.Order{
		ID:           "1",
		Name:         "Order 1",
		FullName:     "test test",
		Address:      "via roma 11",
		Status:       "Created",
		UpdateTime:   time.Now().UnixMilli(),
		CreationTime: time.Now().UnixMilli(),
		Reservations: []string{},
		Warehouses:   []model.OrderWarehouseUsed{},
		Goods: []model.GoodStock{
			{
				GoodID:   "1",
				Quantity: 10,
			},
		},
	}, nil)
	getOrderUseCaseMock.EXPECT().GetAllOrders(gomock.Any()).Return([]model.Order{
		{
			ID:           "1",
			Name:         "Order 1",
			FullName:     "test test",
			Address:      "via roma 11",
			Status:       "Created",
			UpdateTime:   time.Now().UnixMilli(),
			CreationTime: time.Now().UnixMilli(),
			Reservations: []string{},
			Warehouses:   []model.OrderWarehouseUsed{},
			Goods: []model.GoodStock{
				{
					GoodID:   "1",
					Quantity: 10,
				},
			},
		},
		{
			ID:           "2",
			Name:         "Order 2",
			FullName:     "test test",
			Address:      "via roma 11",
			Status:       "Created",
			UpdateTime:   time.Now().UnixMilli(),
			CreationTime: time.Now().UnixMilli(),
			Reservations: []string{},
			Warehouses:   []model.OrderWarehouseUsed{},
			Goods: []model.GoodStock{
				{
					GoodID:   "2",
					Quantity: 10,
				},
			},
		},
	})

	ns, _ := broker.NewInProcessNATSServer(t)

	app := fx.New(
		lib.ModuleTest,
		fx.Supply(ns, t),
		fx.Supply(fx.Annotate(createOrderUseCaseMock, fx.As(new(port.ICreateOrderUseCase)))),
		fx.Supply(fx.Annotate(getOrderUseCaseMock, fx.As(new(port.IGetOrderUseCase)))),
		fx.Provide(NewOrderController),
		fx.Provide(NewOrderRouter),
		fx.Invoke(func(lc fx.Lifecycle, r *OrderRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					require.NoError(t, err)

					dto := request.GetOrderRequestDTO{
						OrderID: "1",
					}
					payload, err := json.Marshal(dto)
					require.NoError(t, err)

					resp, err := ns.Request("order.get", payload, 1*time.Second)
					require.NoError(t, err)

					var respDto response.GetOrderResponseDTO
					err = json.Unmarshal(resp.Data, &respDto)
					require.NoError(t, err)

					require.Empty(t, respDto.Error)
					require.Equal(t, respDto.Message.OrderID, "1")

					resp, err = ns.Request("order.get.all", []byte{}, 1*time.Second)
					require.NoError(t, err)

					var respAllDto response.GetAllOrderResponseDTO
					err = json.Unmarshal(resp.Data, &respAllDto)
					require.NoError(t, err)

					require.Empty(t, respAllDto.Error)
					slices.SortFunc(respAllDto.Message, func(a response.OrderInfo, b response.OrderInfo) int {
						return cmp.Compare(a.OrderID, b.OrderID)
					})
					require.Equal(t, respAllDto.Message[0].OrderID, "1")
					require.Equal(t, respAllDto.Message[1].OrderID, "2")

					return nil
				},
			})
		}),
	)

	err := app.Start(ctx)
	require.NoError(t, err)

	defer func() {
		err = app.Stop(ctx)
		require.NoError(t, err)
	}()
}
