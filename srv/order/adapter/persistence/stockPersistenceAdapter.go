package persistence

import (
	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
)

type StockPersistanceAdapter struct {
	stockRepo IStockRepository
}

func NewStockPersistanceAdapter(stockRepo IStockRepository) *StockPersistanceAdapter {
	return &StockPersistanceAdapter{stockRepo}
}

func (s *StockPersistanceAdapter) ApplyStockUpdate(cmd port.ApplyStockUpdateCmd) error {
	for _, good := range cmd.Goods {
		s.stockRepo.SetStock(cmd.WarehouseID, string(good.ID), good.Quantity)
	}

	return nil
}

func (s *StockPersistanceAdapter) GetStock(cmd port.GetStockCmd) (model.GoodStock, error) {
	stock, err := s.stockRepo.GetStock(cmd.WarehouseID, string(cmd.GoodID))
	if err != nil {
		if err == ErrWarehouseNotFound {
			return model.GoodStock{}, port.ErrStockNotFound
		}
		return model.GoodStock{}, err
	}

	return model.GoodStock{ID: cmd.GoodID, Quantity: stock}, nil
}
