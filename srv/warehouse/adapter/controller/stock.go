package controller

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"github.com/nats-io/nats.go"
)

type StockController struct {
	broker             *broker.NatsMessageBroker
	addStockUseCase    port.AddStockUseCase
	removeStockUseCase port.RemoveStockUseCase
}

func NewStockController(n *broker.NatsMessageBroker, addStockUseCase port.AddStockUseCase, removeStockUseCase port.RemoveStockUseCase) *StockController {
	return &StockController{n, addStockUseCase, removeStockUseCase}
}

func (c *StockController) AddStockHandler(ctx context.Context, msg *nats.Msg) error {
	var dto request.AddStockRequestDTO
	err := json.Unmarshal(msg.Data, &dto)
	if err != nil {
		return err
	}

	cmd := port.AddStockCmd{
		ID:       dto.GoodID,
		Quantity: dto.Quantity,
	}

	err = c.addStockUseCase.AddStock(ctx, cmd)
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

	cmd := port.RemoveStockCmd{
		ID:       dto.GoodID,
		Quantity: dto.Quantity,
	}

	err = c.removeStockUseCase.RemoveStock(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}
