package persistence

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/model"
)

type StockPersistanceAdapter struct {
	stockRepo StockRepository
}

func NewStockPersistanceAdapter(stockRepo StockRepository) *StockPersistanceAdapter {
	return &StockPersistanceAdapter{stockRepo}
}

func (s *StockPersistanceAdapter) ApplyStockUpdate(goods []model.GoodStock) error {
	for _, good := range goods {
		s.stockRepo.SetStock(good.ID, good.Quantity)
	}

	return nil
}

func (s *StockPersistanceAdapter) GetStock(goodId string) int64 {
	return s.stockRepo.GetStock(goodId)
}
