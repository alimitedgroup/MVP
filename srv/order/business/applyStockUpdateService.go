package business

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
)

type ApplyStockUpdateService struct {
	applyStockUpdatePort      port.IApplyStockUpdatePort
	applyOrderUpdatePort      port.IApplyOrderUpdatePort
	getOrderPort              port.IGetOrderPort
	setCompletedWarehousePort port.ISetCompletedWarehouseOrderPort
}

func NewApplyStockUpdateService(
	applyStockUpdatePort port.IApplyStockUpdatePort, getOrderPort port.IGetOrderPort,
	applyOrderUpdatePort port.IApplyOrderUpdatePort, setCompletedWarehousePort port.ISetCompletedWarehouseOrderPort,
) *ApplyStockUpdateService {
	return &ApplyStockUpdateService{applyStockUpdatePort, applyOrderUpdatePort, getOrderPort, setCompletedWarehousePort}
}

func (s *ApplyStockUpdateService) ApplyStockUpdate(ctx context.Context, cmd port.StockUpdateCmd) error {
	// check if the stock update is related to an order
	if cmd.Type == port.StockUpdateCmdTypeOrder {
		order, err := s.getOrderPort.GetOrder(model.OrderID(cmd.OrderID))
		if err != nil {
			return err
		}

		isRelatedToReserv := false
		for _, reserv := range order.Reservations {
			if reserv == cmd.ReservationID {
				isRelatedToReserv = true
				break
			}
		}

		if isRelatedToReserv {
			goods := make([]model.GoodStock, 0, len(cmd.Goods))

			for _, good := range cmd.Goods {
				goods = append(goods, model.GoodStock{
					ID:       model.GoodId(good.GoodID),
					Quantity: good.Delta,
				})
			}

			completedCmd := port.SetCompletedWarehouseCmd{
				WarehouseId: cmd.WarehouseID,
				OrderId:     model.OrderID(cmd.OrderID),
				Goods:       goods,
			}

			order, err := s.setCompletedWarehousePort.SetCompletedWarehouse(completedCmd)
			if err != nil {
				return err
			}

			status := order.Status
			if order.IsCompleted() {
				status = "Completed"
			}

			// stock update for the order reservation
			applyCmd := port.ApplyOrderUpdateCmd{
				Id:           cmd.OrderID,
				Status:       status,
				Name:         order.Name,
				Email:        order.Email,
				Address:      order.Address,
				Goods:        order.Goods,
				Reservations: order.Reservations,
				UpdateTime:   order.UpdateTime,
				CreationTime: order.CreationTime,
			}
			err = s.applyOrderUpdatePort.ApplyOrderUpdate(applyCmd)
			if err != nil {
				return err
			}
		}
	}

	goods := make([]model.GoodStock, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, model.GoodStock{
			ID:       model.GoodId(good.GoodID),
			Quantity: good.Quantity,
		})
	}

	portCmd := port.ApplyStockUpdateCmd{
		WarehouseID: cmd.WarehouseID,
		Goods:       goods,
	}
	err := s.applyStockUpdatePort.ApplyStockUpdate(portCmd)
	if err != nil {
		return err
	}

	return nil
}
