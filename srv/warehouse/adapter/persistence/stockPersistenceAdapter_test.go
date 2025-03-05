package persistence_test

import (
	"testing"

	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/persistence"
	"github.com/alimitedgroup/MVP/srv/warehouse/model"
	"github.com/magiconair/properties/assert"
	"go.uber.org/fx"
)

type stockRepositoryMock struct {
	M map[string]int64
}

func NewStockRepositoryMock() *stockRepositoryMock {
	return &stockRepositoryMock{M: make(map[string]int64)}
}

func (s *stockRepositoryMock) SetStock(id string, quantity int64) bool {
	s.M[id] = quantity
	return true
}

func (s *stockRepositoryMock) GetStock(id string) int64 {
	return s.M[id]
}

func (s *stockRepositoryMock) AddStock(id string, quantity int64) bool {
	s.M[id] += quantity
	return true
}

func TestStockPersistanceAdapter(t *testing.T) {
	ctx := t.Context()

	app := fx.New(
		fx.Provide(NewStockRepositoryMock),
		fx.Provide(func(s *stockRepositoryMock) persistence.StockRepository { return s }),
		fx.Provide(persistence.NewStockPersistanceAdapter),
		fx.Invoke(func(a *persistence.StockPersistanceAdapter, stockRepo persistence.StockRepository, stockRepoMock *stockRepositoryMock) {
			goods := []model.GoodStock{
				{ID: "1", Quantity: 10},
				{ID: "2", Quantity: 20},
			}

			err := a.ApplyStockUpdate(goods)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, stockRepo.GetStock("1"), int64(10))
			assert.Equal(t, stockRepoMock.M["2"], int64(20))
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
