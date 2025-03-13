package business

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/google/uuid"
	"go.uber.org/fx"
)

type ManageOrderService struct {
	getOrderPort        port.IGetOrderPort
	saveOrderUpdatePort port.ISaveOrderUpdatePort
}

type ManageOrderServiceParams struct {
	fx.In
	GetOrderPort        port.IGetOrderPort
	SaveOrderUpdatePort port.ISaveOrderUpdatePort
}

func NewManageStockService(p ManageOrderServiceParams) *ManageOrderService {
	return &ManageOrderService{p.GetOrderPort, p.SaveOrderUpdatePort}
}

func (s *ManageOrderService) CreateOrder(ctx context.Context, cmd port.CreateOrderCmd) (port.CreateOrderResponse, error) {
	orderId := uuid.New().String()

	goods := make([]port.SaveOrderUpdateGood, 0)
	for _, good := range cmd.Goods {
		goods = append(goods, port.SaveOrderUpdateGood{
			GoodId:   good.GoodID,
			Quantity: good.Quantity,
		})
	}

	saveCmd := port.SaveOrderUpdateCmd{
		ID:      orderId,
		Name:    cmd.Name,
		Email:   cmd.Email,
		Address: cmd.Address,
		Goods:   goods,
	}

	err := s.saveOrderUpdatePort.SaveOrderUpdate(ctx, saveCmd)
	if err != nil {
		return port.CreateOrderResponse{}, err
	}

	resp := port.CreateOrderResponse{
		OrderID: orderId,
	}

	return resp, nil
}

func (s *ManageOrderService) GetOrder(ctx context.Context) (model.Order, error) {
	return model.Order{}, nil
}
