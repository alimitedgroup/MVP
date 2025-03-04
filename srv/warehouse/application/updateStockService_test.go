package application_test

import (
	"testing"

	"github.com/alimitedgroup/MVP/srv/warehouse/application"
	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/model"
	"github.com/magiconair/properties/assert"
	"go.uber.org/fx"
)

type saveUpdateStockPortMock struct {
	M     map[string]int64
	Total int64
}

func newSaveUpdateStockPortMock() *saveUpdateStockPortMock {
	return &saveUpdateStockPortMock{M: make(map[string]int64), Total: 0}
}

func (m *saveUpdateStockPortMock) SaveUpdateStock(goods []model.GoodStock) error {
	for _, v := range goods {

		old, exist := m.M[v.ID]
		if !exist {
			old = 0
		}

		m.M[v.ID] = old + v.Quantity
		m.Total += v.Quantity
	}
	return nil
}

func TestUpdateStockService(t *testing.T) {
	ctx := t.Context()

	app := fx.New(
		fx.Provide(newSaveUpdateStockPortMock, func(s *saveUpdateStockPortMock) port.SaveUpdateStockPort { return s }),
		fx.Provide(fx.Annotate(application.NewUpdateStockService, fx.As(new(port.UpdateStockUseCase)))),
		fx.Invoke(func(useCase port.UpdateStockUseCase, saveStockUpdatePortMock *saveUpdateStockPortMock) {
			cmd := port.UpdateStockCmd{
				ID: "1",
				Goods: []port.UpdateStockCommandGood{
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

			err := useCase.UpdateStock(cmd)
			if err != nil {
				t.Errorf("error updating stock: %v", err)
			}

			assert.Equal(t, saveStockUpdatePortMock.Total, int64(30))
			assert.Equal(t, saveStockUpdatePortMock.M["1"], int64(10))
			assert.Equal(t, saveStockUpdatePortMock.M["2"], int64(20))
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
