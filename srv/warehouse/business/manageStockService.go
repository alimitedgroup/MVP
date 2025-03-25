package business

import (
	"context"
	"fmt"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"go.uber.org/fx"
)

type ManageStockService struct {
	createStockUpdatePort port.ICreateStockUpdatePort
	getStockPort          port.IGetStockPort
	getGoodPort           port.IGetGoodPort
}

type ManageStockServiceParams struct {
	fx.In

	CreateStockUpdatePort port.ICreateStockUpdatePort
	GetStockPort          port.IGetStockPort
	GetGoodPort           port.IGetGoodPort
}

func NewManageStockService(p ManageStockServiceParams) *ManageStockService {
	return &ManageStockService{p.CreateStockUpdatePort, p.GetStockPort, p.GetGoodPort}
}

func (s *ManageStockService) AddStock(ctx context.Context, cmd port.AddStockCmd) error {
	if s.getGoodPort.GetGood(model.GoodID(cmd.GoodID)) == nil {
		return fmt.Errorf("good %s not found", cmd.GoodID)
	}

	currentQuantity := s.getStockPort.GetStock(model.GoodID(cmd.GoodID))
	quantity := currentQuantity.Quantity + cmd.Quantity

	createCmd := port.CreateStockUpdateCmd{
		Type: port.CreateStockUpdateCmdTypeAdd,
		Goods: []port.CreateStockUpdateGood{
			{
				Good: model.GoodStock{
					ID:       cmd.GoodID,
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
	if s.getGoodPort.GetGood(model.GoodID(cmd.GoodID)) == nil {
		return fmt.Errorf("good %s not found", cmd.GoodID)
	}

	currentQuantity := s.getStockPort.GetStock(model.GoodID(cmd.GoodID))

	if currentQuantity.Quantity < cmd.Quantity {
		return fmt.Errorf("not enough stock for good %s", cmd.GoodID)
	}

	quantity := currentQuantity.Quantity - cmd.Quantity

	createCmd := port.CreateStockUpdateCmd{
		Type: port.CreateStockUpdateCmdTypeRemove,
		Goods: []port.CreateStockUpdateGood{
			{
				Good: model.GoodStock{
					ID:       cmd.GoodID,
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
