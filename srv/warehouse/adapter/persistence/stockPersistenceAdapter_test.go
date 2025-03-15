package persistence

import (
	"testing"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
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

func (s *stockRepositoryMock) GetFreeStock(goodId string) int64 {
	return 0
}

func (s *stockRepositoryMock) ReserveStock(goodId string, stock int64) error {
	return nil
}

func (s *stockRepositoryMock) UnReserveStock(goodId string, stock int64) error {
	return nil
}

func TestStockPersistanceAdapter(t *testing.T) {
	ctx := t.Context()

	mock := NewStockRepositoryMock()

	app := fx.New(
		fx.Supply(fx.Annotate(mock, fx.As(new(IStockRepository)))),
		fx.Provide(NewStockPersistanceAdapter),
		fx.Invoke(func(a *StockPersistanceAdapter, stockRepo IStockRepository) {
			goods := []model.GoodStock{
				{ID: "1", Quantity: 10},
				{ID: "2", Quantity: 20},
			}

			err := a.ApplyStockUpdate(goods)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, stockRepo.GetStock("1"), int64(10))
			assert.Equal(t, mock.M["2"], int64(20))
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
