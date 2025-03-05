package application

import (
	"context"
	"fmt"

	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/model"
)

type ManageStockService struct {
	createStockUpdatePort port.CreateStockUpdatePort
	getStockPort          port.GetStockPort
}

func NewManageStockService(createStockUpdatePort port.CreateStockUpdatePort, getStockPort port.GetStockPort) *ManageStockService {
	return &ManageStockService{createStockUpdatePort, getStockPort}
}

func (s *ManageStockService) AddStock(ctx context.Context, cmd port.AddStockCmd) error {
	currentQuantity := s.getStockPort.GetStock(cmd.ID)
	quantity := currentQuantity + cmd.Quantity

	createCmd := port.CreateStockUpdateCmd{
		Type: port.CreateStockUpdateCmdTypeAdd,
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

func (s *ManageStockService) RemoveStock(ctx context.Context, cmd port.RemoveStockCmd) error {
	currentQuantity := s.getStockPort.GetStock(cmd.ID)

	if currentQuantity < cmd.Quantity {
		return fmt.Errorf("not enough stock for good %s", cmd.ID)
	}

	quantity := currentQuantity - cmd.Quantity

	createCmd := port.CreateStockUpdateCmd{
		Type: port.CreateStockUpdateCmdTypeRemove,
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
