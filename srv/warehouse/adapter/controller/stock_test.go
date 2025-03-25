package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap/zaptest"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/magiconair/properties/assert"
	"go.uber.org/fx"
)

type mockStock struct {
	stockMap map[string]int64
}

func NewMockStock() *mockStock {
	return &mockStock{stockMap: make(map[string]int64)}
}

func (m *mockStock) AddStock(ctx context.Context, cmd port.AddStockCmd) error {
	old, ok := m.stockMap[cmd.ID]
	if !ok {
		old = 0
	}

	m.stockMap[cmd.ID] = old + cmd.Quantity
	return nil
}

func (m *mockStock) RemoveStock(ctx context.Context, cmd port.RemoveStockCmd) error {
	old, ok := m.stockMap[cmd.ID]
	if !ok {
		return fmt.Errorf("stock not found")
	}

	if old < cmd.Quantity {
		return fmt.Errorf("stock not enough")
	}

	m.stockMap[cmd.ID] = old - cmd.Quantity
	return nil
}

func TestStockController(t *testing.T) {
	ctx := t.Context()

	mock := NewMockStock()

	ns, _ := broker.NewInProcessNATSServer(t)
	cfg := config.WarehouseConfig{ID: "1"}

	app := fx.New(
		fx.Supply(&cfg),
		fx.Supply(ns),
		fx.Supply(zaptest.NewLogger(t)),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(NewStockController),
		fx.Provide(NewStockRouter),
		fx.Supply(fx.Annotate(mock, fx.As(new(port.IAddStockUseCase)), fx.As(new(port.IRemoveStockUseCase)))),
		fx.Invoke(func(lc fx.Lifecycle, r *StockRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}

					addDto := request.AddStockRequestDTO{
						GoodID:   "1",
						Quantity: 10,
					}
					addPayload, err := json.Marshal(addDto)
					if err != nil {
						t.Error(err)
					}

					addResp, err := ns.Request(fmt.Sprintf("warehouse.%s.stock.add", cfg.ID), addPayload, 1*time.Second)
					if err != nil {
						t.Error(err)
					}

					var addRespDto response.ResponseDTO[string]
					err = json.Unmarshal(addResp.Data, &addRespDto)
					if err != nil {
						t.Error(err)
					}

					assert.Equal(t, addRespDto.Message, "ok")

					remDto := request.AddStockRequestDTO{
						GoodID:   "1",
						Quantity: 10,
					}
					remPayload, err := json.Marshal(remDto)
					if err != nil {
						t.Error(err)
					}

					remResp, err := ns.Request(fmt.Sprintf("warehouse.%s.stock.remove", cfg.ID), remPayload, 1*time.Second)
					if err != nil {
						t.Error(err)
					}

					var remRespDto response.ResponseDTO[string]
					err = json.Unmarshal(remResp.Data, &remRespDto)
					if err != nil {
						t.Error(err)
					}

					assert.Equal(t, remRespDto.Message, "ok")

					return nil
				},
				OnStop: func(ctx context.Context) error {
					return nil
				},
			})
		}),
	)

	err := app.Start(ctx)
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
