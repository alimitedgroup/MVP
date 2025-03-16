package business

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
)

type ApplyOrderUpdateService struct {
	applyOrderUpdatePort port.IApplyOrderUpdatePort
}

func NewApplyOrderUpdateService(applyOrderUpdatePort port.IApplyOrderUpdatePort) *ApplyOrderUpdateService {
	return &ApplyOrderUpdateService{applyOrderUpdatePort}
}

func (s *ApplyOrderUpdateService) ApplyOrderUpdate(ctx context.Context, cmd port.OrderUpdateCmd) error {
	portCmd := orderUpdateCmdToApplyOrderUpdateCmd(cmd)
	if err := s.applyOrderUpdatePort.ApplyOrderUpdate(portCmd); err != nil {
		return err
	}

	return nil
}

func orderUpdateCmdToApplyOrderUpdateCmd(cmd port.OrderUpdateCmd) port.ApplyOrderUpdateCmd {
	goods := make([]model.GoodStock, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, model.GoodStock{
			ID:       model.GoodId(good.GoodID),
			Quantity: good.Quantity,
		})
	}

	applyOrderUpdateCmd := port.ApplyOrderUpdateCmd{
		Id:           cmd.ID,
		Status:       cmd.Status,
		Name:         cmd.Name,
		Email:        cmd.Email,
		Address:      cmd.Address,
		CreationTime: cmd.CreationTime,
		Reservations: cmd.Reservations,
		Goods:        goods,
	}

	return applyOrderUpdateCmd
}
