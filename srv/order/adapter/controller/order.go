package controller

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
)

type OrderController struct {
	broker             *broker.NatsMessageBroker
	createOrderUseCase port.ICreateOrderUseCase
	getOrderUseCase    port.IGetOrderUseCase
}

type OrderControllerParams struct {
	fx.In

	Broker             *broker.NatsMessageBroker
	CreateOrderUseCase port.ICreateOrderUseCase
	GetOrderUseCase    port.IGetOrderUseCase
}

func NewOrderController(p OrderControllerParams) *OrderController {
	return &OrderController{p.Broker, p.CreateOrderUseCase, p.GetOrderUseCase}
}

func (c *OrderController) OrderCreateHandler(ctx context.Context, msg *nats.Msg) error {
	var dto request.CreateOrderRequestDTO
	if err := json.Unmarshal(msg.Data, &dto); err != nil {
		return err
	}

	if err := checkCreateOrderRequestDTO(dto); err != nil {
		resp := response.ErrorResponseDTO{
			Error: err.Error(),
		}
		if err := broker.RespondToMsg(msg, resp); err != nil {
			return err
		}

		return nil

	}

	goods := make([]port.CreateOrderGood, 0, len(dto.Goods))
	for _, good := range dto.Goods {
		goods = append(goods, port.CreateOrderGood{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}

	cmd := port.CreateOrderCmd{
		Status:  "Created",
		Name:    dto.Name,
		Email:   dto.Email,
		Address: dto.Address,
		Goods:   goods,
	}

	resp, err := c.createOrderUseCase.CreateOrder(ctx, cmd)
	if err != nil {
		resp := response.ErrorResponseDTO{
			Error: err.Error(),
		}
		if err := broker.RespondToMsg(msg, resp); err != nil {
			return err
		}

		return nil
	}

	respDto := response.OrderCreateResponseDTO{
		Message: response.OrderCreateInfo{
			OrderID: resp.OrderID,
		},
	}

	if err = broker.RespondToMsg(msg, respDto); err != nil {
		return err
	}

	return nil
}

func (c *OrderController) OrderGetHandler(ctx context.Context, msg *nats.Msg) error {
	var dto request.GetOrderRequestDTO
	if err := json.Unmarshal(msg.Data, &dto); err != nil {
		return err
	}

	if dto.OrderID == "" {
		resp := response.ErrorResponseDTO{
			Error: "orderid is required",
		}
		if err := broker.RespondToMsg(msg, resp); err != nil {
			return err
		}
		return nil
	}

	order, err := c.getOrderUseCase.GetOrder(ctx, dto.OrderID)
	if err != nil {
		resp := response.ErrorResponseDTO{
			Error: err.Error(),
		}
		if err := broker.RespondToMsg(msg, resp); err != nil {
			return err
		}
	}

	respDto := orderToGetGoodResponseDTO(order)
	if err := broker.RespondToMsg(msg, respDto); err != nil {
		return err
	}

	return nil
}

var ErrNameIsRequired = errors.New("name is required")
var ErrEmailIsRequired = errors.New("email is required")
var ErrAddressIsRequired = errors.New("address is required")
var ErrGoodsIsRequired = errors.New("goods is required")

func checkCreateOrderRequestDTO(dto request.CreateOrderRequestDTO) error {
	if dto.Name == "" {
		return ErrNameIsRequired
	}

	if dto.Email == "" {
		return ErrEmailIsRequired
	}

	if dto.Address == "" {
		return ErrAddressIsRequired
	}

	if len(dto.Goods) == 0 {
		return ErrGoodsIsRequired
	}

	return nil
}

func orderToGetGoodResponseDTO(order model.Order) response.GetOrderResponseDTO {
	goods := make([]response.OrderInfoGood, 0, len(order.Goods))
	for _, good := range order.Goods {
		goods = append(goods, response.OrderInfoGood{
			GoodID:   string(good.ID),
			Quantity: good.Quantity,
		})
	}

	return response.GetOrderResponseDTO{
		Message: response.OrderInfo{
			OrderID: string(order.Id),
			Status:  order.Status,
			Name:    order.Name,
			Email:   order.Email,
			Address: order.Address,
			Goods:   goods,
		},
	}
}
