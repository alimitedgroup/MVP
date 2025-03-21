package controller

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/nats-io/nats.go"
)

type StockController struct {
	broker             *broker.NatsMessageBroker
	addStockUseCase    port.IAddStockUseCase
	removeStockUseCase port.IRemoveStockUseCase
}

func NewStockController(n *broker.NatsMessageBroker, addStockUseCase port.IAddStockUseCase, removeStockUseCase port.IRemoveStockUseCase) *StockController {
	return &StockController{n, addStockUseCase, removeStockUseCase}
}

func (c *StockController) AddStockHandler(ctx context.Context, msg *nats.Msg) error {
	var dto request.AddStockRequestDTO
	err := json.Unmarshal(msg.Data, &dto)
	if err != nil {
		return err
	}

	cmd := port.AddStockCmd(dto)
	err = c.addStockUseCase.AddStock(ctx, cmd)
	if err != nil {
		resp := response.ResponseDTO[any]{
			Error: err.Error(),
		}

		err = broker.RespondToMsg(msg, resp)
		if err != nil {
			return err
		}

		return nil
	}

	resp := response.ResponseDTO[string]{
		Message: "ok",
	}

	err = broker.RespondToMsg(msg, resp)
	if err != nil {
		return err
	}

	return nil
}

func (c *StockController) RemoveStockHandler(ctx context.Context, msg *nats.Msg) error {
	var dto request.RemoveStockRequestDTO
	err := json.Unmarshal(msg.Data, &dto)
	if err != nil {
		return err
	}

	cmd := port.RemoveStockCmd(dto)
	err = c.removeStockUseCase.RemoveStock(ctx, cmd)
	if err != nil {
		resp := response.ResponseDTO[any]{
			Error: err.Error(),
		}

		err = broker.RespondToMsg(msg, resp)
		if err != nil {
			return err
		}

		return nil
	}

	resp := response.ResponseDTO[string]{
		Message: "ok",
	}

	err = broker.RespondToMsg(msg, resp)
	if err != nil {
		return err
	}

	return nil
}
