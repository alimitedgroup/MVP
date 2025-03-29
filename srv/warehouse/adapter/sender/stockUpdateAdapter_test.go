package sender

import (
	"context"
	"github.com/alimitedgroup/MVP/common/lib"
	"testing"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/magiconair/properties/assert"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestStockUpdateAdapter(t *testing.T) {
	ctx := t.Context()
	cfg := config.WarehouseConfig{
		ID: "1",
	}

	ns, _ := broker.NewInProcessNATSServer(t)
	js, err := jetstream.New(ns)
	if err != nil {
		t.Error(err)
	}

	s, err := js.CreateStream(ctx, stream.StockUpdateStreamConfig)
	if err != nil {
		t.Errorf("failed to create stream: %v", err)
	}

	app := fx.New(
		lib.ModuleTest,
		fx.Supply(ns, t, &cfg),
		fx.Provide(NewPublishStockUpdateAdapter),
		fx.Invoke(func(lc fx.Lifecycle, a *PublishStockUpdateAdapter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					cmd := port.CreateStockUpdateCmd{
						Type: port.CreateStockUpdateCmdTypeAdd,
						Goods: []port.CreateStockUpdateGood{
							{Good: model.GoodStock{ID: "1", Quantity: 10}, QuantityDiff: 10},
						},
					}

					err := a.CreateStockUpdate(ctx, cmd)
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

func TestStockUpdateAdapterNetworkErr(t *testing.T) {
	ctx := t.Context()
	cfg := &config.WarehouseConfig{
		ID: "1",
	}

	ns, _ := broker.NewInProcessNATSServer(t)

	broker := broker.NewTest(t, ns)

	ns.Close()

	a := NewPublishStockUpdateAdapter(broker, cfg)

	cmd := port.CreateStockUpdateCmd{
		Type: port.CreateStockUpdateCmdTypeAdd,
		Goods: []port.CreateStockUpdateGood{
			{Good: model.GoodStock{ID: "1", Quantity: 10}, QuantityDiff: 10},
		},
	}

	err := a.CreateStockUpdate(ctx, cmd)
	require.Error(t, err)
}
