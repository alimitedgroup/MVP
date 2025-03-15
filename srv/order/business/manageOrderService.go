package business

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/google/uuid"
	"go.uber.org/fx"
)

type ManageOrderService struct {
	getOrderPort                 port.IGetOrderPort
	saveOrderUpdatePort          port.ISaveOrderUpdatePort
	calculateAvailabilityUseCase port.ICalculateAvailabilityUseCase
}

type ManageOrderServiceParams struct {
	fx.In
	GetOrderPort                 port.IGetOrderPort
	SaveOrderUpdatePort          port.ISaveOrderUpdatePort
	CalculateAvailabilityUseCase port.ICalculateAvailabilityUseCase
}

func NewManageStockService(p ManageOrderServiceParams) *ManageOrderService {
	return &ManageOrderService{p.GetOrderPort, p.SaveOrderUpdatePort, p.CalculateAvailabilityUseCase}
}

func (s *ManageOrderService) CreateOrder(ctx context.Context, cmd port.CreateOrderCmd) (port.CreateOrderResponse, error) {
	availCmd := createOrderCmdToCalculateAvailabilityCmd(cmd)
	availResp, err := s.calculateAvailabilityUseCase.GetAvailable(ctx, availCmd)
	if err != nil {
		return port.CreateOrderResponse{}, err
	}

	_ = availResp

	orderId := uuid.New().String()
	saveCmd := createOrderCmdToSaveOrderUpdateCmd(orderId, cmd)

	err = s.saveOrderUpdatePort.SaveOrderUpdate(ctx, saveCmd)
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

func (s *ManageOrderService) GetAllOrders(context.Context) ([]model.Order, error) {
	orders, err := s.getOrderPort.GetAllOrder()
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func createOrderCmdToCalculateAvailabilityCmd(cmd port.CreateOrderCmd) port.CalculateAvailabilityCmd {
	requestGoods := make([]port.RequestedGood, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		requestGoods = append(requestGoods, port.RequestedGood(good))
	}

	availCmd := port.CalculateAvailabilityCmd{
		Goods: requestGoods,
	}

	return availCmd
}

func createOrderCmdToSaveOrderUpdateCmd(orderId string, cmd port.CreateOrderCmd) port.SaveOrderUpdateCmd {
	goods := make([]port.SaveOrderUpdateGood, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, port.SaveOrderUpdateGood{
			GoodId:   good.GoodID,
			Quantity: good.Quantity,
		})
	}

	saveCmd := port.SaveOrderUpdateCmd{
		ID:      orderId,
		Status:  cmd.Status,
		Name:    cmd.Name,
		Email:   cmd.Email,
		Address: cmd.Address,
		Goods:   goods,
	}

	return saveCmd
}
