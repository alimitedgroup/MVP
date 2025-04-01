package business

import (
	"context"
	"log"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"go.uber.org/fx"
)

type ApplyStockUpdateService struct {
	applyStockUpdatePort      port.IApplyStockUpdatePort
	applyOrderUpdatePort      port.IApplyOrderUpdatePort
	getOrderPort              port.IGetOrderPort
	getTransferPort           port.IGetTransferPort
	applyTransferUpdatePort   port.IApplyTransferUpdatePort
	setCompleteTransferPort   port.ISetCompleteTransferPort
	setCompletedWarehousePort port.ISetCompletedWarehouseOrderPort
}

type ApplyStockUpdateServiceParams struct {
	fx.In

	ApplyStockUpdatePort      port.IApplyStockUpdatePort
	ApplyOrderUpdatePort      port.IApplyOrderUpdatePort
	GetOrderPort              port.IGetOrderPort
	GetTransferPort           port.IGetTransferPort
	ApplyTransferUpdatePort   port.IApplyTransferUpdatePort
	SetCompleteTransferPort   port.ISetCompleteTransferPort
	SetCompletedWarehousePort port.ISetCompletedWarehouseOrderPort
}

func NewApplyStockUpdateService(
	p ApplyStockUpdateServiceParams,
) *ApplyStockUpdateService {
	return &ApplyStockUpdateService{p.ApplyStockUpdatePort, p.ApplyOrderUpdatePort, p.GetOrderPort, p.GetTransferPort, p.ApplyTransferUpdatePort, p.SetCompleteTransferPort, p.SetCompletedWarehousePort}
}

func (s *ApplyStockUpdateService) ApplyStockUpdate(ctx context.Context, cmd port.StockUpdateCmd) error {
	// check if the stock update is related to an order
	if cmd.Type == port.StockUpdateCmdTypeOrder {
		err := s.applyStockUpdateFromOrder(cmd)
		if err != nil {
			return err
		}
	} else if cmd.Type == port.StockUpdateCmdTypeTransfer {
		err := s.applyStockUpdateFromTransfer(cmd)
		if err != nil {
			return err
		}
	}

	goods := make([]model.GoodStock, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, model.GoodStock{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}

	portCmd := port.ApplyStockUpdateCmd{
		WarehouseID: cmd.WarehouseID,
		Goods:       goods,
	}
	s.applyStockUpdatePort.ApplyStockUpdate(portCmd)

	return nil
}

func (s *ApplyStockUpdateService) applyStockUpdateFromTransfer(cmd port.StockUpdateCmd) error {
	transfer, err := s.getTransferPort.GetTransfer(model.TransferID(cmd.TransferID))
	if err != nil {
		return err
	}

	// if transfer.ReservationID == cmd.ReservationID {

	// }
	if err := s.setCompleteTransferPort.IncrementLinkedStockUpdate(model.TransferID(transfer.ID)); err != nil {
		return err
	}
	transfer, err = s.getTransferPort.GetTransfer(model.TransferID(transfer.ID))
	if err != nil {
		return err
	}

	if transfer.LinkedStockUpdate == 2 {
		if err := s.setCompleteTransferPort.SetComplete(model.TransferID(transfer.ID)); err != nil {
			return err
		}
	}
	return nil
}

func (s *ApplyStockUpdateService) applyStockUpdateFromOrder(cmd port.StockUpdateCmd) error {
	order, err := s.getOrderPort.GetOrder(model.OrderID(cmd.OrderID))
	if err != nil {
		return err
	}
	log.Printf("order in applystock: %v\n", order)

	goods := make([]model.GoodStock, 0, len(cmd.Goods))

	for _, good := range cmd.Goods {
		goods = append(goods, model.GoodStock{
			GoodID:   good.GoodID,
			Quantity: good.Delta,
		})
	}

	completedCmd := port.SetCompletedWarehouseCmd{
		WarehouseID: cmd.WarehouseID,
		OrderID:     cmd.OrderID,
		Goods:       goods,
	}
	order, err = s.setCompletedWarehousePort.SetCompletedWarehouse(completedCmd)
	if err != nil {
		return err
	}
	log.Printf("order in applystock after set: %v\n", order)

	completed := order.IsCompleted()
	log.Printf("is completed: %v\n", completed)
	if completed {
		if err := s.setCompletedWarehousePort.SetComplete(model.OrderID(cmd.OrderID)); err != nil {
			return err
		}
		order, err = s.getOrderPort.GetOrder(model.OrderID(cmd.OrderID))
		if err != nil {
			return err
		}
		if order.Status != "Completed" {
			log.Printf("error didn't set order to completed: %v\n", order)
		}
	}

	return nil
}
