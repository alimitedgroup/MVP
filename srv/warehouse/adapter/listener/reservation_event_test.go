package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
)

func TestReservationEventListener(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)
	mock := NewMockIApplyReservationUseCase(ctrl)
	mock.EXPECT().ApplyReservationEvent(gomock.Any()).Return(nil)

	ns, _ := broker.NewInProcessNATSServer(t)
	js, err := jetstream.New(ns)
	require.NoError(t, err)

	cfg := config.WarehouseConfig{ID: "1"}

	app := fx.New(
		fx.Supply(&cfg),
		fx.Supply(ns),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(fx.Annotate(broker.NewRestoreStreamControlFactory, fx.As(new(broker.IRestoreStreamControlFactory)))),
		fx.Supply(fx.Annotate(mock, fx.As(new(port.IApplyReservationUseCase)))),
		fx.Provide(NewReservationEventListener),
		fx.Provide(NewReservationEventRouter),
		fx.Invoke(func(lc fx.Lifecycle, r *ReservationEventRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					require.NoError(t, err)

					event := stream.ReservationEvent{
						ID: "1",
						Goods: []stream.ReservationGood{
							{
								GoodID:   "1",
								Quantity: 1,
							},
						},
					}

					payload, err := json.Marshal(event)
					require.NoError(t, err)

					ack, err := js.Publish(ctx, fmt.Sprintf("reservation.%s", cfg.ID), payload)
					require.NoError(t, err)

					time.Sleep(100 * time.Millisecond)

					require.Equal(t, ack.Stream, "reservation")

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
