package business

import (
	"context"
	"log"
	"time"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/google/uuid"
	"go.uber.org/fx"
)

type ManageOrderService struct {
	getOrderPort                 port.IGetOrderPort
	sendOrderUpdatePort          port.ISendOrderUpdatePort
	sendContactWarehousePort     port.ISendContactWarehousePort
	requestReservationPort       port.IRequestReservationPort
	calculateAvailabilityUseCase port.ICalculateAvailabilityUseCase
}

type ManageOrderServiceParams struct {
	fx.In
	GetOrderPort                 port.IGetOrderPort
	SendOrderUpdatePort          port.ISendOrderUpdatePort
	SendContactWarehousePort     port.ISendContactWarehousePort
	RequestReservationPort       port.IRequestReservationPort
	CalculateAvailabilityUseCase port.ICalculateAvailabilityUseCase
}

func NewManageStockService(p ManageOrderServiceParams) *ManageOrderService {
	return &ManageOrderService{p.GetOrderPort, p.SendOrderUpdatePort, p.SendContactWarehousePort, p.RequestReservationPort, p.CalculateAvailabilityUseCase}
}

func (s *ManageOrderService) CreateOrder(ctx context.Context, cmd port.CreateOrderCmd) (port.CreateOrderResponse, error) {
	orderId := uuid.New().String()
	saveCmd := createOrderCmdToSendOrderUpdateCmd(orderId, cmd)

	order, err := s.sendOrderUpdatePort.SendOrderUpdate(ctx, saveCmd)
	if err != nil {
		return port.CreateOrderResponse{}, err
	}

	contactCmd := port.SendContactWarehouseCmd{
		Order:                 order,
		TransferId:            "",
		LastContact:           0,
		ConfirmedReservations: []port.ConfirmedReservation{},
		ExcludeWarehouses:     []string{},
	}
	err = s.sendContactWarehousePort.SendContactWarehouses(ctx, contactCmd)
	if err != nil {
		return port.CreateOrderResponse{}, err
	}

	resp := port.CreateOrderResponse{
		OrderID: orderId,
	}

	return resp, nil
}

func (s *ManageOrderService) GetOrder(ctx context.Context, orderId string) (model.Order, error) {
	order, err := s.getOrderPort.GetOrder(model.OrderID(orderId))
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (s *ManageOrderService) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	orders, err := s.getOrderPort.GetAllOrder()
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *ManageOrderService) ContactWarehouses(ctx context.Context, cmd port.ContactWarehousesCmd) error {
	now := time.Now().UnixMilli()

	availCmd := contactCmdToCalculateAvailabilityCmd(cmd)
	availResp, err := s.calculateAvailabilityUseCase.GetAvailable(ctx, availCmd)
	if err != nil {
		orderUpdateCmd := contactCmdToSendOrderUpdateCmdForCancel(cmd)
		_, err := s.sendOrderUpdatePort.SendOrderUpdate(ctx, orderUpdateCmd)
		if err != nil {
			return err
		}
		return nil
	}

	log.Printf("availResp: %v", availResp)

	confirmed := make([]port.ConfirmedReservation, 0, len(availResp.Warehouses)+len(cmd.ConfirmedReservations))
	confirmed = append(confirmed, cmd.ConfirmedReservations...)
	remainingConfirmation := len(availResp.Warehouses)

	errWarehouses := make([]string, 0, len(cmd.ExcludeWarehouses))

	for _, warehouse := range availResp.Warehouses {
		reservCmd := warehouseAvailabilityToReservationCmd(warehouse)
		reservResp, err := s.requestReservationPort.RequestReservation(ctx, reservCmd)
		if err != nil {
			// this reservation gave an error continue to the next warehouse
			errWarehouses = append(errWarehouses, warehouse.WarehouseID)
			continue
		}

		confirmed = append(confirmed, port.ConfirmedReservation{
			WarehouseId:   warehouse.WarehouseID,
			ReservationID: reservResp.Id,
			Goods:         warehouse.Goods,
		})

		remainingConfirmation--
	}

	if remainingConfirmation > 0 {
		// not all goods are reserved
		// send another contact request, to take a snaptshot of the current reservations and retry later
		goods := make([]model.GoodStock, 0, len(cmd.Order.Goods))
		for _, good := range cmd.Order.Goods {
			goods = append(goods, model.GoodStock{
				ID:       model.GoodId(good.GoodId),
				Quantity: good.Quantity,
			})
		}

		sendContactCmd := port.SendContactWarehouseCmd{
			Order: model.Order{
				Id:           model.OrderID(cmd.Order.ID),
				Status:       cmd.Order.Status,
				Name:         cmd.Order.Name,
				FullName:     cmd.Order.FullName,
				Address:      cmd.Order.Address,
				UpdateTime:   cmd.Order.UpdateTime,
				CreationTime: cmd.Order.CreationTime,
				Goods:        goods,
				Reservations: cmd.Order.Reservations,
				Warehouses:   []model.OrderWarehouseUsed{},
			},
			LastContact:           now,
			ConfirmedReservations: confirmed,
			ExcludeWarehouses:     errWarehouses,
		}
		if err = s.sendContactWarehousePort.SendContactWarehouses(ctx, sendContactCmd); err != nil {
			return err
		}
	} else {
		// all goods are reserved
		orderUpdateCmd := contactCmdAndConfirmedToSendOrderUpdateCmd(cmd, confirmed)
		_, err := s.sendOrderUpdatePort.SendOrderUpdate(ctx, orderUpdateCmd)
		if err != nil {
			return err
		}
	}

	return nil
}

func createOrderCmdToSendOrderUpdateCmd(orderId string, cmd port.CreateOrderCmd) port.SendOrderUpdateCmd {
	goods := make([]port.SendOrderUpdateGood, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, port.SendOrderUpdateGood{
			GoodId:   good.GoodID,
			Quantity: good.Quantity,
		})
	}

	saveCmd := port.SendOrderUpdateCmd{
		ID:           orderId,
		Status:       "Created",
		Name:         cmd.Name,
		FullName:     cmd.FullName,
		Address:      cmd.Address,
		Goods:        goods,
		Reservations: []string{},
	}

	return saveCmd
}

func contactCmdToCalculateAvailabilityCmd(cmd port.ContactWarehousesCmd) port.CalculateAvailabilityCmd {
	used := make(map[string]int64)
	for _, reserv := range cmd.ConfirmedReservations {
		for good, quantity := range reserv.Goods {
			if _, ok := used[good]; !ok {
				used[good] = 0
			}
			used[good] += quantity
		}
	}

	goods := make([]port.RequestedGood, 0, len(cmd.Order.Goods))
	for _, good := range cmd.Order.Goods {
		goodReserved, ok := used[good.GoodId]
		if !ok {
			goodReserved = 0
		}
		if goodReserved >= good.Quantity {
			continue
		}

		goods = append(goods, port.RequestedGood{
			GoodID:   good.GoodId,
			Quantity: good.Quantity - goodReserved,
		})
	}

	availCmd := port.CalculateAvailabilityCmd{Goods: goods, ExcludedWarehouses: cmd.ExcludeWarehouses}
	return availCmd
}

func warehouseAvailabilityToReservationCmd(warehouse port.WarehouseAvailability) port.RequestReservationCmd {
	items := make([]port.ReservationItem, 0, len(warehouse.Goods))
	for good, stock := range warehouse.Goods {
		items = append(items, port.ReservationItem{
			GoodId:   good,
			Quantity: stock,
		})
	}
	reservCmd := port.RequestReservationCmd{
		WarehouseId: warehouse.WarehouseID,
		Items:       items,
	}
	return reservCmd
}

func contactCmdAndConfirmedToSendOrderUpdateCmd(cmd port.ContactWarehousesCmd, confirmed []port.ConfirmedReservation) port.SendOrderUpdateCmd {
	goods := make([]port.SendOrderUpdateGood, 0, len(cmd.Order.Goods))
	for _, good := range cmd.Order.Goods {
		goods = append(goods, port.SendOrderUpdateGood(good))
	}

	reservations := make([]string, 0, len(confirmed))
	for _, reserv := range confirmed {
		reservations = append(reservations, reserv.ReservationID)
	}

	orderUpdatecmd := port.SendOrderUpdateCmd{
		ID:           cmd.Order.ID,
		Status:       "Filled",
		Name:         cmd.Order.Name,
		FullName:     cmd.Order.FullName,
		CreationTime: cmd.Order.CreationTime,
		Goods:        goods,
		Reservations: reservations,
	}
	return orderUpdatecmd
}

func contactCmdToSendOrderUpdateCmdForCancel(cmd port.ContactWarehousesCmd) port.SendOrderUpdateCmd {
	goods := make([]port.SendOrderUpdateGood, 0, len(cmd.Order.Goods))
	for _, good := range cmd.Order.Goods {
		goods = append(goods, port.SendOrderUpdateGood(good))
	}

	orderUpdatecmd := port.SendOrderUpdateCmd{
		ID:           cmd.Order.ID,
		Status:       "Cancelled",
		Name:         cmd.Order.Name,
		FullName:     cmd.Order.FullName,
		CreationTime: cmd.Order.CreationTime,
		Goods:        goods,
		Reservations: cmd.Order.Reservations,
	}
	return orderUpdatecmd
}
