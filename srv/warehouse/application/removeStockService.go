package application

import (
	"context"
	"fmt"

	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/model"
)

type RemoveStockService struct {
	createStockUpdatePort port.CreateStockUpdatePort
	getStockPort          port.GetStockPort
}

func NewRemoveStockService(createStockUpdatePort port.CreateStockUpdatePort, getStockPort port.GetStockPort) *RemoveStockService {
	return &RemoveStockService{createStockUpdatePort, getStockPort}
}

func (s *RemoveStockService) RemoveStock(ctx context.Context, cmd port.RemoveStockCmd) error {
	currentQuantity := s.getStockPort.GetStock(cmd.ID)

	if currentQuantity < cmd.Quantity {
		return fmt.Errorf("not enough stock for good %s", cmd.ID)
	}

	quantity := currentQuantity - cmd.Quantity

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
