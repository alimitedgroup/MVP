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
	getGoodPort           port.GetGoodPort
}

func NewManageStockService(createStockUpdatePort port.CreateStockUpdatePort, getStockPort port.GetStockPort, getGoodPort port.GetGoodPort) *ManageStockService {
	return &ManageStockService{createStockUpdatePort, getStockPort, getGoodPort}
}

func (s *ManageStockService) AddStock(ctx context.Context, cmd port.AddStockCmd) error {
	if s.getGoodPort.GetGood(cmd.ID) == nil {
		return fmt.Errorf("good %s not found", cmd.ID)
	}

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
	if s.getGoodPort.GetGood(cmd.ID) == nil {
		return fmt.Errorf("good %s not found", cmd.ID)
	}

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
