package business

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
)

type ApplyStockUpdateService struct {
	applyStockUpdatePort port.IApplyStockUpdatePort
}

func NewApplyStockUpdateService(applyStockUpdatePort port.IApplyStockUpdatePort) *ApplyStockUpdateService {
	return &ApplyStockUpdateService{applyStockUpdatePort}
}

func (s *ApplyStockUpdateService) ApplyStockUpdate(ctx context.Context, cmd port.StockUpdateCmd) error {
	goods := make([]model.GoodStock, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, model.GoodStock{
			ID:       model.GoodId(good.GoodID),
			Quantity: good.Quantity,
		})
	}

	portCmd := port.ApplyStockUpdateCmd{
		WarehouseID: cmd.WarehouseID,
		Goods:       goods,
	}

	err := s.applyStockUpdatePort.ApplyStockUpdate(portCmd)
	if err != nil {
		return err
	}

	return nil
}
