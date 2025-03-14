package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

type applyStockUpdateMock struct {
	sync.Mutex
	stockMap map[string]int64
}

func NewApplyStockUpdateMock() *applyStockUpdateMock {
	return &applyStockUpdateMock{stockMap: make(map[string]int64)}
}

func (m *applyStockUpdateMock) ApplyStockUpdate(_ context.Context, cmd port.StockUpdateCmd) error {
	m.Lock()
	defer m.Unlock()

	for _, good := range cmd.Goods {
		old := m.stockMap[good.GoodID]

		if cmd.Type == port.StockUpdateCmdTypeAdd {
			m.stockMap[good.GoodID] = old + good.Quantity
		} else if cmd.Type == port.StockUpdateCmdTypeRemove {
			m.stockMap[good.GoodID] = old - good.Quantity
		}
	}
	return nil
}

func (m *applyStockUpdateMock) GetStock(id string) int64 {
	m.Lock()
	defer m.Unlock()

	stock := m.stockMap[id]
	return stock
}

func TestStockUpdateListener(t *testing.T) {
	ctx := t.Context()

	ns, _ := broker.NewInProcessNATSServer(t)
	js, err := jetstream.New(ns)
	if err != nil {
		t.Error(err)
	}

	cfg := config.WarehouseConfig{ID: "1"}
	mock := NewApplyStockUpdateMock()

	app := fx.New(
		fx.Supply(&cfg),
		fx.Supply(ns),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(fx.Annotate(broker.NewRestoreStreamControlFactory, fx.As(new(broker.IRestoreStreamControlFactory)))),
		fx.Supply(fx.Annotate(mock, fx.As(new(port.IApplyStockUpdateUseCase)))),
		fx.Provide(NewStockUpdateListener),
		fx.Provide(NewStockUpdateRouter),
		fx.Invoke(func(lc fx.Lifecycle, r *StockUpdateRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}

					event := stream.StockUpdate{
						Type:        stream.StockUpdateTypeAdd,
						ID:          "1",
						WarehouseID: cfg.ID,
						Goods: []stream.StockUpdateGood{
							{
								GoodID:   "1",
								Quantity: 10,
							},
							{
								GoodID:   "2",
								Quantity: 20,
							},
						},
					}

					payload, err := json.Marshal(event)
					if err != nil {
						t.Error(err)
					}

					ack, err := js.Publish(ctx, fmt.Sprintf("stock.update.%s", cfg.ID), payload)
					if err != nil {
						t.Error(err)
					}

					// for {
					// 	info, err := r.StockUpdateConsumer.Info(ctx)
					// 	if err != nil {
					// 		t.Error(err)
					// 	}
					// 	if info.NumAckPending == 0 {
					// 		break
					// 	}
					// 	time.Sleep(10 * time.Millisecond)
					// }
					time.Sleep(100 * time.Millisecond)

					assert.Equal(t, ack.Stream, "stock_update")
					assert.Equal(t, int64(10), mock.GetStock("1"))
					assert.Equal(t, int64(20), mock.GetStock("2"))

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
