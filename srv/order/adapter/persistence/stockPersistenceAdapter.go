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

func (s *StockPersistanceAdapter) ApplyStockUpdate(cmd port.ApplyStockUpdateCmd) {
	for _, good := range cmd.Goods {
		s.stockRepo.SetStock(cmd.WarehouseID, good.GoodID, good.Quantity)
	}
}

func (s *StockPersistanceAdapter) GetStock(cmd port.GetStockCmd) (model.GoodStock, error) {
	stock, err := s.stockRepo.GetStock(cmd.WarehouseID, cmd.GoodID)
	if err != nil {
		if err == ErrWarehouseNotFound {
			return model.GoodStock{}, port.ErrStockNotFound
		}
		return model.GoodStock{}, err
	}

	return model.GoodStock{GoodID: cmd.GoodID, Quantity: stock}, nil
}

func (s *StockPersistanceAdapter) GetGlobalStock(GoodID model.GoodID) model.GoodStock {
	stock := s.stockRepo.GetGlobalStock(string(GoodID))
	return model.GoodStock{GoodID: string(GoodID), Quantity: stock}
}

func (s *StockPersistanceAdapter) GetWarehouses() []model.Warehouse {
	warehousesIds := s.stockRepo.GetWarehouses()

	warehouses := make([]model.Warehouse, 0, len(warehousesIds))
	for _, warehouseId := range warehousesIds {
		warehouses = append(warehouses, model.Warehouse{ID: warehouseId})
	}

	return warehouses
}
