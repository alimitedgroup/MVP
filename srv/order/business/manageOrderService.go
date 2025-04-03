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
	getTransferPort              port.IGetTransferPort
	sendOrderUpdatePort          port.ISendOrderUpdatePort
	sendTransferUpdatePort       port.ISendTransferUpdatePort
	sendContactWarehousePort     port.ISendContactWarehousePort
	requestReservationPort       port.IRequestReservationPort
	calculateAvailabilityUseCase port.ICalculateAvailabilityUseCase
	transactionPort              port.TransactionPort
}

type ManageOrderServiceParams struct {
	fx.In
	GetOrderPort                 port.IGetOrderPort
	GetTransferPort              port.IGetTransferPort
	SendOrderUpdatePort          port.ISendOrderUpdatePort
	SendTransferUpdatePort       port.ISendTransferUpdatePort
	SendContactWarehousePort     port.ISendContactWarehousePort
	RequestReservationPort       port.IRequestReservationPort
	CalculateAvailabilityUseCase port.ICalculateAvailabilityUseCase
	TransactionPort              port.TransactionPort
}

func NewManageOrderService(p ManageOrderServiceParams) *ManageOrderService {
	return &ManageOrderService{p.GetOrderPort, p.GetTransferPort, p.SendOrderUpdatePort, p.SendTransferUpdatePort, p.SendContactWarehousePort, p.RequestReservationPort, p.CalculateAvailabilityUseCase, p.TransactionPort}
}

func (s *ManageOrderService) CreateTransfer(ctx context.Context, cmd port.CreateTransferCmd) (port.CreateTransferResponse, error) {
	s.transactionPort.Lock()
	defer s.transactionPort.Unlock()

	transferId := uuid.New().String()

	goods := make([]port.SendTransferUpdateGood, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, port.SendTransferUpdateGood(good))
	}

	transferCmd := port.SendTransferUpdateCmd{
		ID:            transferId,
		Status:        "Created",
		SenderID:      cmd.SenderID,
		ReceiverID:    cmd.ReceiverID,
		ReservationID: "",
		Goods:         goods,
	}
	transfer, err := s.sendTransferUpdatePort.SendTransferUpdate(ctx, transferCmd)
	if err != nil {
		return port.CreateTransferResponse{}, err
	}

	contactCmd := port.SendContactWarehouseCmd{
		Type:                  port.SendContactWarehouseTypeTransfer,
		Order:                 nil,
		Transfer:              &transfer,
		ConfirmedReservations: []port.ConfirmedReservation{},
		ExcludeWarehouses:     []string{},
		RetryInTime:           0,
		RetryUntil:            time.Now().Add(time.Hour * 12).UnixMilli(),
	}
	err = s.sendContactWarehousePort.SendContactWarehouses(ctx, contactCmd)
	if err != nil {
		return port.CreateTransferResponse{}, err
	}

	return port.CreateTransferResponse{TransferID: transferId}, nil
}

func (s *ManageOrderService) CreateOrder(ctx context.Context, cmd port.CreateOrderCmd) (port.CreateOrderResponse, error) {
	s.transactionPort.Lock()
	defer s.transactionPort.Unlock()

	orderId := uuid.New().String()
	saveCmd := createOrderCmdToSendOrderUpdateCmd(orderId, cmd)

	order, err := s.sendOrderUpdatePort.SendOrderUpdate(ctx, saveCmd)
	if err != nil {
		return port.CreateOrderResponse{}, err
	}

	contactCmd := port.SendContactWarehouseCmd{
		Type:                  port.SendContactWarehouseTypeOrder,
		Order:                 &order,
		Transfer:              nil,
		ConfirmedReservations: []port.ConfirmedReservation{},
		ExcludeWarehouses:     []string{},
		RetryInTime:           0,
		RetryUntil:            time.Now().Add(time.Hour * 12).UnixMilli(),
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

func (s *ManageOrderService) GetOrder(ctx context.Context, orderId port.GetOrderCmd) (model.Order, error) {
	s.transactionPort.Lock()
	defer s.transactionPort.Unlock()

	order, err := s.getOrderPort.GetOrder(model.OrderID(orderId))
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (s *ManageOrderService) GetAllOrders(ctx context.Context) []model.Order {
	s.transactionPort.Lock()
	defer s.transactionPort.Unlock()

	orders := s.getOrderPort.GetAllOrder()
	return orders
}

func (s *ManageOrderService) GetTransfer(ctx context.Context, transferId port.GetTransferCmd) (model.Transfer, error) {
	s.transactionPort.Lock()
	defer s.transactionPort.Unlock()

	transfer, err := s.getTransferPort.GetTransfer(model.TransferID(transferId))
	if err != nil {
		return model.Transfer{}, err
	}

	return transfer, nil
}

func (s *ManageOrderService) GetAllTransfers(ctx context.Context) []model.Transfer {
	s.transactionPort.Lock()
	defer s.transactionPort.Unlock()

	transfers := s.getTransferPort.GetAllTransfer()
	return transfers
}

func (s *ManageOrderService) ContactWarehouses(ctx context.Context, cmd port.ContactWarehousesCmd) (port.ContactWarehousesResponse, error) {
	s.transactionPort.Lock()
	defer s.transactionPort.Unlock()

	var err error
	var resp port.ContactWarehousesResponse

	if time.UnixMilli(cmd.RetryInTime).After(time.Now()) {
		retryAfter := time.Until(time.UnixMilli(cmd.RetryInTime))
		return port.ContactWarehousesResponse{IsRetry: true, RetryAfter: retryAfter}, nil
	}

	switch cmd.Type {
	case port.ContactWarehousesTypeOrder:
		resp, err = s.contactWarehouseForOrder(ctx, cmd)
	case port.ContactWarehousesTypeTransfer:
		resp, err = s.contactWarehouseForTransfer(ctx, cmd)
	}

	if err != nil {
		return port.ContactWarehousesResponse{}, err
	}
	return resp, nil
}

func (s *ManageOrderService) contactWarehouseForTransfer(ctx context.Context, cmd port.ContactWarehousesCmd) (port.ContactWarehousesResponse, error) {
	items := make([]port.ReservationGood, 0, len(cmd.Transfer.Goods))
	for _, good := range cmd.Transfer.Goods {
		items = append(items, port.ReservationGood(good))
	}
	reservCmd := port.RequestReservationCmd{
		WarehouseID: cmd.Transfer.SenderID,
		Goods:       items,
	}
	reservResp, err := s.requestReservationPort.RequestReservation(ctx, reservCmd)
	if err != nil {
		now := time.Now()
		if err == port.ErrNotEnoughStock || time.UnixMilli(cmd.RetryUntil).Before(now) {
			// TODO: error message in the order
			sendTransferCmd := contactCmdToSendTransferUpdateCmdForCancel(cmd)
			_, err = s.sendTransferUpdatePort.SendTransferUpdate(ctx, sendTransferCmd)
			if err != nil {
				return port.ContactWarehousesResponse{}, err
			}
			return port.ContactWarehousesResponse{}, nil
		}

		goods := make([]model.GoodStock, 0, len(cmd.Transfer.Goods))
		for _, good := range cmd.Transfer.Goods {
			goods = append(goods, model.GoodStock{
				GoodID:   good.GoodID,
				Quantity: good.Quantity,
			})
		}
		sendContactCmd := port.SendContactWarehouseCmd{
			Type: port.SendContactWarehouseTypeTransfer,
			Transfer: &model.Transfer{
				ID:            cmd.Transfer.ID,
				Status:        cmd.Transfer.Status,
				UpdateTime:    cmd.Transfer.UpdateTime,
				CreationTime:  cmd.Transfer.CreationTime,
				SenderID:      cmd.Transfer.SenderID,
				ReceiverID:    cmd.Transfer.ReceiverID,
				Goods:         goods,
				ReservationID: "",
			},
			Order:                 nil,
			RetryInTime:           now.Add(time.Second * 10).UnixMilli(),
			RetryUntil:            cmd.RetryUntil,
			ConfirmedReservations: []port.ConfirmedReservation{},
			ExcludeWarehouses:     []string{},
		}
		if err = s.sendContactWarehousePort.SendContactWarehouses(ctx, sendContactCmd); err != nil {
			return port.ContactWarehousesResponse{}, err
		}
		return port.ContactWarehousesResponse{}, nil
	}

	goods := make([]port.SendTransferUpdateGood, 0, len(cmd.Transfer.Goods))
	for _, good := range cmd.Transfer.Goods {
		goods = append(goods, port.SendTransferUpdateGood(good))
	}
	sendTransferCmd := port.SendTransferUpdateCmd{
		Status:        "Filled",
		ID:            cmd.Transfer.ID,
		CreationTime:  cmd.Transfer.CreationTime,
		SenderID:      cmd.Transfer.SenderID,
		ReceiverID:    cmd.Transfer.ReceiverID,
		Goods:         goods,
		ReservationID: reservResp.ID,
	}
	_, err = s.sendTransferUpdatePort.SendTransferUpdate(ctx, sendTransferCmd)
	if err != nil {
		return port.ContactWarehousesResponse{}, err
	}

	return port.ContactWarehousesResponse{IsRetry: false}, nil
}

func (s *ManageOrderService) contactWarehouseForOrder(ctx context.Context, cmd port.ContactWarehousesCmd) (port.ContactWarehousesResponse, error) {

	availCmd := contactCmdToCalculateAvailabilityCmd(cmd)
	availResp, err := s.calculateAvailabilityUseCase.GetAvailable(ctx, availCmd)
	if err != nil {
		orderUpdateCmd := contactCmdToSendOrderUpdateCmdForCancel(cmd)
		_, err := s.sendOrderUpdatePort.SendOrderUpdate(ctx, orderUpdateCmd)
		if err != nil {
			return port.ContactWarehousesResponse{}, err
		}
		return port.ContactWarehousesResponse{}, nil
	}

	confirmed := make([]port.ConfirmedReservation, 0, len(availResp.Warehouses)+len(cmd.ConfirmedReservations))
	confirmed = append(confirmed, cmd.ConfirmedReservations...)
	remainingConfirmation := len(availResp.Warehouses)

	errWarehouses := make([]string, 0, len(cmd.ExcludeWarehouses))

	for _, warehouse := range availResp.Warehouses {
		log.Printf("Warehouse %s items available: %v\n", warehouse.WarehouseID, warehouse.Goods)
		reservCmd := warehouseAvailabilityToReservationCmd(warehouse)
		reservResp, err := s.requestReservationPort.RequestReservation(ctx, reservCmd)
		if err != nil {
			// this reservation gave an error continue to the next warehouse
			errWarehouses = append(errWarehouses, warehouse.WarehouseID)
			continue
		}

		confirmed = append(confirmed, port.ConfirmedReservation{
			WarehouseID:   warehouse.WarehouseID,
			ReservationID: reservResp.ID,
			Goods:         warehouse.Goods,
		})

		remainingConfirmation--
	}

	if remainingConfirmation > 0 {
		now := time.Now()
		// not all goods are reserved
		// send another contact request, to take a snaptshot of the current reservations and retry later
		if time.UnixMilli(cmd.RetryUntil).Before(now) {
			orderUpdateCmd := contactCmdToSendOrderUpdateCmdForCancel(cmd)
			_, err := s.sendOrderUpdatePort.SendOrderUpdate(ctx, orderUpdateCmd)
			if err != nil {
				return port.ContactWarehousesResponse{}, err
			}
			return port.ContactWarehousesResponse{}, nil
		}

		goods := make([]model.GoodStock, 0, len(cmd.Order.Goods))
		for _, good := range cmd.Order.Goods {
			goods = append(goods, model.GoodStock{
				GoodID:   good.GoodID,
				Quantity: good.Quantity,
			})
		}

		sendContactCmd := port.SendContactWarehouseCmd{
			Type: port.SendContactWarehouseTypeOrder,
			Order: &model.Order{
				ID:           cmd.Order.ID,
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
			ConfirmedReservations: confirmed,
			ExcludeWarehouses:     errWarehouses,
			RetryInTime:           now.Add(time.Second * 10).UnixMilli(),
			RetryUntil:            cmd.RetryUntil,
		}
		if err = s.sendContactWarehousePort.SendContactWarehouses(ctx, sendContactCmd); err != nil {
			return port.ContactWarehousesResponse{}, err
		}
		return port.ContactWarehousesResponse{}, nil
	} else {
		// all goods are reserved
		orderUpdateCmd := contactCmdAndConfirmedToSendOrderUpdateCmd(cmd, confirmed)
		_, err := s.sendOrderUpdatePort.SendOrderUpdate(ctx, orderUpdateCmd)
		if err != nil {
			return port.ContactWarehousesResponse{}, err
		}
	}

	return port.ContactWarehousesResponse{}, nil

}

func createOrderCmdToSendOrderUpdateCmd(orderId string, cmd port.CreateOrderCmd) port.SendOrderUpdateCmd {
	goods := make([]port.SendOrderUpdateGood, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, port.SendOrderUpdateGood(good))
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
		goodReserved, ok := used[good.GoodID]
		if !ok {
			goodReserved = 0
		}
		if goodReserved >= good.Quantity {
			continue
		}

		goods = append(goods, port.RequestedGood{
			GoodID:   good.GoodID,
			Quantity: good.Quantity - goodReserved,
		})
	}

	availCmd := port.CalculateAvailabilityCmd{Goods: goods, ExcludedWarehouses: cmd.ExcludeWarehouses}
	return availCmd
}

func warehouseAvailabilityToReservationCmd(warehouse port.WarehouseAvailability) port.RequestReservationCmd {
	items := make([]port.ReservationGood, 0, len(warehouse.Goods))
	for good, stock := range warehouse.Goods {
		items = append(items, port.ReservationGood{
			GoodID:   good,
			Quantity: stock,
		})
	}
	reservCmd := port.RequestReservationCmd{
		WarehouseID: warehouse.WarehouseID,
		Goods:       items,
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
		Address:      cmd.Order.Address,
		Name:         cmd.Order.Name,
		FullName:     cmd.Order.FullName,
		CreationTime: cmd.Order.CreationTime,
		UpdateTime:   cmd.Order.UpdateTime,
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
		Address:      cmd.Order.Address,
		Name:         cmd.Order.Name,
		FullName:     cmd.Order.FullName,
		CreationTime: cmd.Order.CreationTime,
		UpdateTime:   cmd.Order.UpdateTime,
		Goods:        goods,
		Reservations: cmd.Order.Reservations,
	}
	return orderUpdatecmd
}

func contactCmdToSendTransferUpdateCmdForCancel(cmd port.ContactWarehousesCmd) port.SendTransferUpdateCmd {
	goods := make([]port.SendTransferUpdateGood, 0, len(cmd.Transfer.Goods))
	for _, good := range cmd.Transfer.Goods {
		goods = append(goods, port.SendTransferUpdateGood(good))
	}
	sendTransferCmd := port.SendTransferUpdateCmd{
		Status:        "Cancelled",
		ID:            cmd.Transfer.ID,
		CreationTime:  cmd.Transfer.CreationTime,
		UpdateTime:    cmd.Transfer.UpdateTime,
		SenderID:      cmd.Transfer.SenderID,
		ReceiverID:    cmd.Transfer.ReceiverID,
		Goods:         goods,
		ReservationID: cmd.Transfer.ReservationID,
	}
	return sendTransferCmd
}
