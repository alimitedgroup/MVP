package persistence

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
)

type StockPersistanceAdapter struct {
	stockRepo IStockRepository
}

func NewStockPersistanceAdapter(stockRepo IStockRepository) *StockPersistanceAdapter {
	return &StockPersistanceAdapter{stockRepo}
}

func (s *StockPersistanceAdapter) ApplyStockUpdate(goods []model.GoodStock) error {
	for _, good := range goods {
		s.stockRepo.SetStock(string(good.ID), good.Quantity)
	}

	return nil
}

func (s *StockPersistanceAdapter) GetStock(goodId model.GoodId) model.GoodStock {
	stock := s.stockRepo.GetStock(string(goodId))
	return model.GoodStock{
		ID:       goodId,
		Quantity: stock,
	}
}
