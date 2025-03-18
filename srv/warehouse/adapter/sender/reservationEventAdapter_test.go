package sender

import (
	"context"
	"testing"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/magiconair/properties/assert"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestReservationEventAdapter(t *testing.T) {
	ctx := t.Context()
	cfg := &config.WarehouseConfig{
		ID: "1",
	}

	ns, _ := broker.NewInProcessNATSServer(t)
	js, err := jetstream.New(ns)
	if err != nil {
		t.Error(err)
	}

	s, err := js.CreateStream(ctx, stream.ReservationEventStreamConfig)
	if err != nil {
		t.Errorf("failed to create stream: %v", err)
	}

	app := fx.New(
		fx.Supply(ns),
		fx.Supply(cfg),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(NewPublishReservationEventAdapter),
		fx.Invoke(func(lc fx.Lifecycle, a *PublishReservationEventAdapter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					cmd := model.Reservation{
						ID: "1",
						Goods: []model.ReservationGood{
							{
								GoodID:   "1",
								Quantity: 10,
							},
						},
					}

					err := a.StoreReservationEvent(ctx, cmd)
					if err != nil {
						t.Error(err)
					}

					info, err := s.Info(ctx)
					if err != nil {
						t.Error(err)
					}
					assert.Equal(t, info.State.Msgs, uint64(1))

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

func TestReservationEventAdapterNetworkErr(t *testing.T) {
	ctx := t.Context()
	cfg := &config.WarehouseConfig{
		ID: "1",
	}

	ns, _ := broker.NewInProcessNATSServer(t)

	broker, err := broker.NewNatsMessageBroker(ns)
	require.Nil(t, err)

	ns.Close()

	a := NewPublishReservationEventAdapter(broker, cfg)

	cmd := model.Reservation{
		ID: "1",
		Goods: []model.ReservationGood{
			{
				GoodID:   "1",
				Quantity: 10,
			},
		},
	}

	err = a.StoreReservationEvent(ctx, cmd)
	require.NotNil(t, err)
}
