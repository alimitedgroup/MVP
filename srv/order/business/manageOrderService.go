package business

import (
	"context"
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

	err := s.sendOrderUpdatePort.SendOrderUpdate(ctx, saveCmd)
	if err != nil {
		return port.CreateOrderResponse{}, err
	}

	contactCmd := port.SendContactWarehouseCmd{
		OrderId:               orderId,
		TransferId:            "",
		LastContact:           0,
		ConfirmedReservations: []port.ConfirmedReservation{},
		ExcludeWarehouses:     []string{},
	}
	s.sendContactWarehousePort.SendContactWarehouses(ctx, contactCmd)

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
	now := time.Now().UnixNano()

	order, err := s.getOrderPort.GetOrder(model.OrderID(cmd.OrderId))
	if err != nil {
		return err
	}

	used := make(map[string]int64)
	for _, reserv := range cmd.ConfirmedReservations {
		for good, quantity := range reserv.Goods {
			if _, ok := used[good]; !ok {
				used[good] = 0
			}
			used[good] += quantity
		}
	}

	goods := make([]port.RequestedGood, 0, len(order.Goods))
	for _, good := range order.Goods {
		goodReserved, ok := used[string(good.ID)]
		if !ok {
			goodReserved = 0
		}
		if goodReserved >= good.Quantity {
			continue
		}

		goods = append(goods, port.RequestedGood{
			GoodID:   string(good.ID),
			Quantity: good.Quantity - goodReserved,
		})
	}

	availCmd := port.CalculateAvailabilityCmd{Goods: goods, ExcludedWarehouses: cmd.ExcludeWarehouses}
	availResp, err := s.calculateAvailabilityUseCase.GetAvailable(ctx, availCmd)
	if err != nil {
		return err
	}

	confirmed := make([]port.ConfirmedReservation, 0, len(availResp.Warehouses)+len(cmd.ConfirmedReservations))
	confirmed = append(confirmed, cmd.ConfirmedReservations...)
	remainingConfirmation := len(availResp.Warehouses)

	errWarehouses := make([]string, 0, len(cmd.ExcludeWarehouses))

	for _, warehouse := range availResp.Warehouses {
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
		sendContactCmd := port.SendContactWarehouseCmd{
			OrderId:               cmd.OrderId,
			LastContact:           now,
			ConfirmedReservations: confirmed,
			ExcludeWarehouses:     errWarehouses,
		}
		if err = s.sendContactWarehousePort.SendContactWarehouses(ctx, sendContactCmd); err != nil {
			return err
		}
	} else {
		// all goods are reserved
		goods := make([]port.SendOrderUpdateGood, 0, len(order.Goods))
		for _, good := range order.Goods {
			goods = append(goods, port.SendOrderUpdateGood{
				GoodId:   string(good.ID),
				Quantity: good.Quantity,
			})
		}

		reservations := make([]string, 0, len(confirmed))
		for _, reserv := range confirmed {
			reservations = append(reservations, reserv.ReservationID)
		}

		orderUpdatecmd := port.SendOrderUpdateCmd{
			ID:           string(order.Id),
			Status:       "Filled",
			Name:         order.Name,
			Email:        order.Email,
			CreationTime: order.CreationTime,
			Goods:        goods,
			Reservations: reservations,
		}
		err := s.sendOrderUpdatePort.SendOrderUpdate(ctx, orderUpdatecmd)
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
		Email:        cmd.Email,
		Address:      cmd.Address,
		Goods:        goods,
		Reservations: []string{},
	}

	return saveCmd
}
