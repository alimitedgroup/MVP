package business

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"go.uber.org/fx"
)

type ApplyOrderUpdateService struct {
	applyOrderUpdatePort    port.IApplyOrderUpdatePort
	applyTransferUpdatePort port.IApplyTransferUpdatePort
	transactionPort         port.ITransactionPort
}

type ApplyOrderUpdateServiceParams struct {
	fx.In

	ApplyOrderUpdatePort    port.IApplyOrderUpdatePort
	ApplyTransferUpdatePort port.IApplyTransferUpdatePort
	TransactionPort         port.ITransactionPort
}

func NewApplyOrderUpdateService(p ApplyOrderUpdateServiceParams) *ApplyOrderUpdateService {
	return &ApplyOrderUpdateService{p.ApplyOrderUpdatePort, p.ApplyTransferUpdatePort, p.TransactionPort}
}

func (s *ApplyOrderUpdateService) ApplyOrderUpdate(ctx context.Context, cmd port.OrderUpdateCmd) {
	s.transactionPort.Lock()
	defer s.transactionPort.Unlock()

	portCmd := orderUpdateCmdToApplyOrderUpdateCmd(cmd)
	s.applyOrderUpdatePort.ApplyOrderUpdate(portCmd)
}

func (s *ApplyOrderUpdateService) ApplyTransferUpdate(ctx context.Context, cmd port.TransferUpdateCmd) {
	s.transactionPort.Lock()
	defer s.transactionPort.Unlock()

	portCmd := transferUpdateCmdToApplyTransferUpdateCmd(cmd)
	s.applyTransferUpdatePort.ApplyTransferUpdate(portCmd)
}

func orderUpdateCmdToApplyOrderUpdateCmd(cmd port.OrderUpdateCmd) port.ApplyOrderUpdateCmd {
	goods := make([]model.GoodStock, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, model.GoodStock{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}

	applyOrderUpdateCmd := port.ApplyOrderUpdateCmd{
		ID:           cmd.ID,
		Status:       cmd.Status,
		Name:         cmd.Name,
		FullName:     cmd.FullName,
		Address:      cmd.Address,
		CreationTime: cmd.CreationTime,
		UpdateTime:   cmd.UpdateTime,
		Reservations: cmd.Reservations,
		Goods:        goods,
	}

	return applyOrderUpdateCmd
}

func transferUpdateCmdToApplyTransferUpdateCmd(cmd port.TransferUpdateCmd) port.ApplyTransferUpdateCmd {
	goods := make([]model.GoodStock, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, model.GoodStock{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}

	applyOrderUpdateCmd := port.ApplyTransferUpdateCmd{
		ID:            cmd.ID,
		Status:        cmd.Status,
		SenderID:      cmd.SenderID,
		ReceiverID:    cmd.ReceiverID,
		UpdateTime:    cmd.UpdateTime,
		CreationTime:  cmd.CreationTime,
		ReservationID: cmd.ReservationID,
		Goods:         goods,
	}

	return applyOrderUpdateCmd
}
