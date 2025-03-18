package business

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
)

type ApplyOrderUpdateService struct {
	applyOrderUpdatePort    port.IApplyOrderUpdatePort
	applyTransferUpdatePort port.IApplyTransferUpdatePort
}

func NewApplyOrderUpdateService(applyOrderUpdatePort port.IApplyOrderUpdatePort, applyTransferUpdatePort port.IApplyTransferUpdatePort) *ApplyOrderUpdateService {
	return &ApplyOrderUpdateService{applyOrderUpdatePort, applyTransferUpdatePort}
}

func (s *ApplyOrderUpdateService) ApplyOrderUpdate(ctx context.Context, cmd port.OrderUpdateCmd) error {
	portCmd := orderUpdateCmdToApplyOrderUpdateCmd(cmd)
	if err := s.applyOrderUpdatePort.ApplyOrderUpdate(portCmd); err != nil {
		return err
	}

	return nil
}

func (s *ApplyOrderUpdateService) ApplyTransferUpdate(ctx context.Context, cmd port.TransferUpdateCmd) error {
	portCmd := transferUpdateCmdToApplyTransferUpdateCmd(cmd)
	if err := s.applyTransferUpdatePort.ApplyTransferUpdate(portCmd); err != nil {
		return err
	}

	return nil
}

func orderUpdateCmdToApplyOrderUpdateCmd(cmd port.OrderUpdateCmd) port.ApplyOrderUpdateCmd {
	goods := make([]model.GoodStock, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, model.GoodStock{
			ID:       model.GoodID(good.GoodID),
			Quantity: good.Quantity,
		})
	}

	applyOrderUpdateCmd := port.ApplyOrderUpdateCmd{
		Id:           cmd.ID,
		Status:       cmd.Status,
		Name:         cmd.Name,
		FullName:     cmd.FullName,
		Address:      cmd.Address,
		CreationTime: cmd.CreationTime,
		Reservations: cmd.Reservations,
		Goods:        goods,
	}

	return applyOrderUpdateCmd
}

func transferUpdateCmdToApplyTransferUpdateCmd(cmd port.TransferUpdateCmd) port.ApplyTransferUpdateCmd {
	goods := make([]model.GoodStock, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, model.GoodStock{
			ID:       model.GoodID(good.GoodID),
			Quantity: good.Quantity,
		})
	}

	applyOrderUpdateCmd := port.ApplyTransferUpdateCmd{
		Id:            cmd.ID,
		Status:        cmd.Status,
		SenderId:      cmd.SenderId,
		ReceiverId:    cmd.ReceiverId,
		UpdateTime:    cmd.UpdateTime,
		CreationTime:  cmd.CreationTime,
		ReservationId: cmd.ReservationId,
		Goods:         goods,
	}

	return applyOrderUpdateCmd
}
