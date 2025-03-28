package controller

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	RemoveStockCounter     metric.Int64Counter
	AddStockRequestCounter metric.Int64Counter
	TotalRequestsCounter   metric.Int64Counter
	Logger                 *zap.Logger
	metricMap              sync.Map
)

type MetricParams struct {
	fx.In
	Logger *zap.Logger
	Meter  metric.Meter
}

type StockController struct {
	addStockUseCase    port.IAddStockUseCase
	removeStockUseCase port.IRemoveStockUseCase
}

func NewStockController(addStockUseCase port.IAddStockUseCase, removeStockUseCase port.IRemoveStockUseCase, mp MetricParams) *StockController {
	observability.CounterSetup(&mp.Meter, mp.Logger, &TotalRequestsCounter, &metricMap, "num_warehouse_requests")
	observability.CounterSetup(&mp.Meter, mp.Logger, &AddStockRequestCounter, &metricMap, "num_add_stock_requests")
	observability.CounterSetup(&mp.Meter, mp.Logger, &RemoveStockCounter, &metricMap, "num_remove_stock_requests")
	Logger = mp.Logger
	return &StockController{addStockUseCase, removeStockUseCase}
}

func (c *StockController) AddStockHandler(ctx context.Context, msg *nats.Msg) error {

	Logger.Info("Received add stock request")
	verdict := "success"

	defer func() {
		Logger.Info("Add stock request terminated")
		AddStockRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		TotalRequestsCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	var dto request.AddStockRequestDTO
	err := json.Unmarshal(msg.Data, &dto)
	if err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		return err
	}

	cmd := port.AddStockCmd(dto)
	err = c.addStockUseCase.AddStock(ctx, cmd)
	if err != nil {
		Logger.Debug("Cannot add stock", zap.Error(err))
		verdict = "cannot add stock"
		resp := response.ResponseDTO[any]{
			Error: err.Error(),
		}

		err = broker.RespondToMsg(msg, resp)
		if err != nil {
			Logger.Error("Cannot send response", zap.Error(err))
			return err
		}

		return nil
	}

	resp := response.ResponseDTO[string]{
		Message: "ok",
	}

	err = broker.RespondToMsg(msg, resp)
	if err != nil {
		Logger.Error("Cannot send response", zap.Error(err))
		return err
	}

	return nil
}

func (c *StockController) RemoveStockHandler(ctx context.Context, msg *nats.Msg) error {

	Logger.Info("Received remove stock request")
	verdict := "success"

	defer func() {
		Logger.Info("Remove stock request terminated")
		AddStockRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		TotalRequestsCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	var dto request.RemoveStockRequestDTO
	err := json.Unmarshal(msg.Data, &dto)
	if err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		return err
	}

	cmd := port.RemoveStockCmd(dto)
	err = c.removeStockUseCase.RemoveStock(ctx, cmd)
	if err != nil {
		Logger.Debug("Cannot remove stock", zap.Error(err))
		verdict = "cannot remove stock"
		resp := response.ResponseDTO[any]{
			Error: err.Error(),
		}

		err = broker.RespondToMsg(msg, resp)
		if err != nil {
			Logger.Error("Cannot send response", zap.Error(err))
			return err
		}

		return nil
	}

	resp := response.ResponseDTO[string]{
		Message: "ok",
	}

	err = broker.RespondToMsg(msg, resp)
	if err != nil {
		Logger.Error("Cannot send response", zap.Error(err))
		return err
	}

	return nil
}
