package business

import (
	"testing"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/magiconair/properties/assert"
	"go.uber.org/fx"
)

type applyStockUpdatePortMock struct {
	M     map[string]int64
	Total int64
}

func newApplyStockUpdatePortMock() *applyStockUpdatePortMock {
	return &applyStockUpdatePortMock{M: make(map[string]int64), Total: 0}
}

func (m *applyStockUpdatePortMock) SaveEventID(port.IdempotentCmd) {

}

func (m *applyStockUpdatePortMock) IsAlreadyProcessed(port.IdempotentCmd) bool {
	return false
}

func (m *applyStockUpdatePortMock) ApplyStockUpdate(goods []model.GoodStock) {
	for _, v := range goods {

		old, exist := m.M[string(v.ID)]
		if !exist {
			old = 0
		}

		m.M[string(v.ID)] = old + v.Quantity
		m.Total += v.Quantity
	}
}

func TestApplyStockUpdateService(t *testing.T) {
	ctx := t.Context()

	mock := newApplyStockUpdatePortMock()

	app := fx.New(
		fx.Supply(fx.Annotate(mock, fx.As(new(port.IApplyStockUpdatePort)), fx.As(new(port.IIdempotentPort)))),
		fx.Provide(fx.Annotate(NewApplyStockUpdateService, fx.As(new(port.IApplyStockUpdateUseCase)))),
		fx.Invoke(func(useCase port.IApplyStockUpdateUseCase) {
			cmd := port.StockUpdateCmd{
				ID: "1",
				Goods: []port.StockUpdateCmdGood{
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

			useCase.ApplyStockUpdate(cmd)

			assert.Equal(t, mock.Total, int64(30))
			assert.Equal(t, mock.M["1"], int64(10))
			assert.Equal(t, mock.M["2"], int64(20))
		}),
	)

	err := app.Start(ctx)
	if err != nil {
		t.Errorf("error starting app: %v", err)
	}

	defer func() {
		err := app.Stop(ctx)
		if err != nil {
			t.Errorf("error stopping app: %v", err)
		}
	}()

}
