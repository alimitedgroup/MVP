package application

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/model"
)

type ApplyStockUpdateService struct {
	applyStockUpdatePort port.ApplyStockUpdatePort
}

func NewApplyStockUpdateService(applyStockUpdatePort port.ApplyStockUpdatePort) *ApplyStockUpdateService {
	return &ApplyStockUpdateService{applyStockUpdatePort}
}

func (s *ApplyStockUpdateService) ApplyStockUpdate(ctx context.Context, cmd port.StockUpdateCmd) error {
	goods := make([]model.GoodStock, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, model.GoodStock{
			ID:       good.GoodID,
			Quantity: good.Quantity,
		})
	}

	err := s.applyStockUpdatePort.ApplyStockUpdate(goods)
	if err != nil {
		return err
	}

	return nil
}
