package listener

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/lib"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
)

func TestOrderUpdateListenerForOrder(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)

	confirmOrderMock := NewMockIConfirmOrderUseCase(ctrl)
	confirmOrderMock.EXPECT().ConfirmOrder(gomock.Any(), gomock.Any()).Return(nil)

	confirmTransferMock := NewMockIConfirmTransferUseCase(ctrl)

	ns, _ := broker.NewInProcessNATSServer(t)
	js, err := jetstream.New(ns)
	require.NoError(t, err)

	cfg := config.WarehouseConfig{ID: "1"}

	app := fx.New(
		lib.ModuleTest,
		fx.Supply(ns, t, &cfg, ctrl),
		fx.Supply(fx.Annotate(confirmOrderMock, fx.As(new(port.IConfirmOrderUseCase)))),
		fx.Supply(fx.Annotate(confirmTransferMock, fx.As(new(port.IConfirmTransferUseCase)))),
		fx.Provide(NewOrderUpdateListener),
		fx.Provide(NewOrderUpdateRouter),
		fx.Provide(observability.TestMeter),
		fx.Provide(observability.TestLogger),
		fx.Invoke(func(lc fx.Lifecycle, r *OrderUpdateRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					require.NoError(t, err)

					event := stream.OrderUpdate{
						ID:           "1",
						Status:       "Created",
						Name:         "Test",
						FullName:     "Test Test",
						Address:      "Test",
						Reservations: []string{},
						UpdateTime:   time.Now().UnixMilli(),
						CreationTime: time.Now().UnixMilli(),
						Goods: []stream.OrderUpdateGood{
							{
								GoodID:   "1",
								Quantity: 1,
							},
						},
					}

					payload, err := json.Marshal(event)
					require.NoError(t, err)

					ack, err := js.Publish(ctx, "order.update", payload)
					require.NoError(t, err)

					time.Sleep(100 * time.Millisecond)

					require.Equal(t, ack.Stream, "order_update")

					return nil
				},
			})
		}),
	)

	err = app.Start(ctx)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = app.Stop(ctx)
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestOrderUpdateListenerForTransfer(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)

	confirmOrderMock := NewMockIConfirmOrderUseCase(ctrl)

	confirmTransferMock := NewMockIConfirmTransferUseCase(ctrl)
	confirmTransferMock.EXPECT().ConfirmTransfer(gomock.Any(), gomock.Any()).Return(nil)

	ns, _ := broker.NewInProcessNATSServer(t)
	js, err := jetstream.New(ns)
	require.NoError(t, err)

	cfg := config.WarehouseConfig{ID: "1"}

	app := fx.New(
		lib.ModuleTest,
		fx.Supply(ns, t, &cfg, ctrl),
		fx.Supply(fx.Annotate(confirmOrderMock, fx.As(new(port.IConfirmOrderUseCase)))),
		fx.Supply(fx.Annotate(confirmTransferMock, fx.As(new(port.IConfirmTransferUseCase)))),
		fx.Provide(NewOrderUpdateListener),
		fx.Provide(NewOrderUpdateRouter),
		fx.Provide(observability.TestMeter),
		fx.Provide(observability.TestLogger),
		fx.Invoke(func(lc fx.Lifecycle, r *OrderUpdateRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					require.NoError(t, err)

					event := stream.TransferUpdate{
						ID:            "1",
						Status:        "Created",
						ReservationID: "",
						SenderID:      "1",
						ReceiverID:    "2",
						UpdateTime:    time.Now().UnixMilli(),
						CreationTime:  time.Now().UnixMilli(),
						Goods: []stream.TransferUpdateGood{
							{
								GoodID:   "1",
								Quantity: 1,
							},
						},
					}

					payload, err := json.Marshal(event)
					require.NoError(t, err)

					ack, err := js.Publish(ctx, "transfer.update", payload)
					require.NoError(t, err)

					time.Sleep(100 * time.Millisecond)

					require.Equal(t, ack.Stream, "transfer_update")

					return nil
				},
			})
		}),
	)

	err = app.Start(ctx)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = app.Stop(ctx)
		if err != nil {
			t.Error(err)
		}
	}()
}
