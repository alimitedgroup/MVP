package business

import (
	"context"
	"fmt"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
)

type ManageStockService struct {
	createStockUpdatePort port.ICreateStockUpdatePort
	getStockPort          port.IGetStockPort
	getGoodPort           port.IGetGoodPort
}

func NewManageStockService(createStockUpdatePort port.ICreateStockUpdatePort, getStockPort port.IGetStockPort, getGoodPort port.IGetGoodPort) *ManageStockService {
	return &ManageStockService{createStockUpdatePort, getStockPort, getGoodPort}
}

func (s *ManageStockService) AddStock(ctx context.Context, cmd port.AddStockCmd) error {
	if s.getGoodPort.GetGood(model.GoodId(cmd.ID)) == nil {
		return fmt.Errorf("good %s not found", cmd.ID)
	}

	currentQuantity := s.getStockPort.GetStock(model.GoodId(cmd.ID))
	quantity := currentQuantity.Quantity + cmd.Quantity

	createCmd := port.CreateStockUpdateCmd{
		Type: port.CreateStockUpdateCmdTypeAdd,
		Goods: []port.CreateStockUpdateCmdGood{
			{
				Good: model.GoodStock{
					ID:       model.GoodId(cmd.ID),
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
	if s.getGoodPort.GetGood(model.GoodId(cmd.ID)) == nil {
		return fmt.Errorf("good %s not found", cmd.ID)
	}

	currentQuantity := s.getStockPort.GetStock(model.GoodId(cmd.ID))

	if currentQuantity.Quantity < cmd.Quantity {
		return fmt.Errorf("not enough stock for good %s", cmd.ID)
	}

	quantity := currentQuantity.Quantity - cmd.Quantity

	createCmd := port.CreateStockUpdateCmd{
		Type: port.CreateStockUpdateCmdTypeRemove,
		Goods: []port.CreateStockUpdateCmdGood{
			{
				Good: model.GoodStock{
					ID:       model.GoodId(cmd.ID),
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
