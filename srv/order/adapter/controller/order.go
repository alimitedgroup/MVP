package controller

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/nats-io/nats.go"
)

type OrderController struct {
	broker             *broker.NatsMessageBroker
	createOrderUseCase port.ICreateOrderUseCase
}

func NewOrderController(broker *broker.NatsMessageBroker, createOrderUseCase port.ICreateOrderUseCase) *OrderController {
	return &OrderController{broker, createOrderUseCase}
}

func (c *OrderController) OrderCreateHandler(ctx context.Context, msg *nats.Msg) error {
	var dto request.CreateOrderRequestDTO
	err := json.Unmarshal(msg.Data, &dto)
	if err != nil {
		return err
	}

	if err := checkCreateOrderRequestDTO(dto); err != nil {
		resp := response.ErrorResponseDTO{
			Error: err.Error(),
		}

		err = broker.RespondToMsg(msg, resp)
		if err != nil {
			return err
		}

		return nil

	}

	goods := make([]port.CreateOrderGood, 0)
	for _, good := range dto.Goods {
		goods = append(goods, port.CreateOrderGood{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}

	cmd := port.CreateOrderCmd{
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

		err = broker.RespondToMsg(msg, resp)
		if err != nil {
			return err
		}

		return nil
	}

	respDto := response.OrderCreateResponseDTO{
		Message: response.OrderCreateInfo{
			OrderID: resp.OrderID,
		},
	}

	err = broker.RespondToMsg(msg, respDto)
	if err != nil {
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
