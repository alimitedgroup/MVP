package application

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/model"
)

type AddStockService struct {
	createStockUpdatePort port.CreateStockUpdatePort
	getStockPort          port.GetStockPort
}

func NewAddStockService(createStockUpdatePort port.CreateStockUpdatePort, getStockPort port.GetStockPort) *AddStockService {
	return &AddStockService{createStockUpdatePort, getStockPort}
}

func (s *AddStockService) AddStock(ctx context.Context, cmd port.AddStockCmd) error {
	currentQuantity := s.getStockPort.GetStock(cmd.ID)
	quantity := currentQuantity + cmd.Quantity

	createCmd := port.CreateStockUpdateCmd{
		Goods: []port.CreateStockUpdateCmdGood{
			{
				Good: model.GoodStock{
					ID:       cmd.ID,
					Quantity: quantity,
				},
				QuantityDiff: cmd.Quantity,
			},
		},
	}

	err := s.createStockUpdatePort.CreateStockUpdate(ctx, createCmd)
	if err != nil {
		return err
	}

	return nil
}
