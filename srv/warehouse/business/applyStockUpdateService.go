package business

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
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

	err := s.applyStockUpdatePort.ApplyStockUpdate(goods)
	if err != nil {
		return err
	}

	return nil
}
