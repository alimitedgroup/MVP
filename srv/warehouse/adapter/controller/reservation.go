package controller

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

var (
	CreateReservationCounter metric.Int64Counter
)

type ReservationController struct {
	createReservationUseCase port.ICreateReservationUseCase
}

func NewReservationController(createReservationUseCase port.ICreateReservationUseCase, mp MetricParams) *ReservationController {
	observability.CounterSetup(&mp.Meter, mp.Logger, &TotalRequestsCounter, &MetricMap, "num_warehouse_requests")
	observability.CounterSetup(&mp.Meter, mp.Logger, &CreateReservationCounter, &MetricMap, "num_create_reservation_requests")
	Logger = mp.Logger
	return &ReservationController{createReservationUseCase}
}

func (c *ReservationController) CreateReservationHandler(ctx context.Context, msg *nats.Msg) error {

	Logger.Info("Received create reservation request")
	verdict := "success"

	defer func() {
		Logger.Info("Create reservation request terminated")
		CreateReservationCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		TotalRequestsCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	var dto request.ReserveStockRequestDTO
	if err := json.Unmarshal(msg.Data, &dto); err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		return err
	}

	for _, good := range dto.Goods {
		if good.Quantity <= 0 {
			verdict = "bad arguments"
			Logger.Debug("quantity must be greater than 0")
			resp := response.ResponseDTO[any]{
				Error: "quantity must be greater than 0",
			}
			return broker.RespondToMsg(msg, resp)
		}
	}

	Logger.Debug("Create reservation request", zap.Any("dto", dto))

	goods := make([]port.ReservationGood, 0, len(dto.Goods))
	for _, good := range dto.Goods {
		goods = append(goods, port.ReservationGood(good))
	}

	cmd := port.CreateReservationCmd{Goods: goods}
	createResp, err := c.createReservationUseCase.CreateReservation(ctx, cmd)
	if err != nil {
		Logger.Debug("Cannot create reservation", zap.Error(err))
		verdict = "cannot create reservation"
		resp := response.ErrorResponseDTO{
			Error: err.Error(),
		}
		if err := broker.RespondToMsg(msg, resp); err != nil {
			Logger.Error("Cannot send response", zap.Error(err))
			return err
		}
	}

	respDto := response.ReserveStockResponseDTO{
		Message: response.ReserveStockInfo{
			ReservationID: createResp.ReservationID,
		},
	}
	if err := broker.RespondToMsg(msg, &respDto); err != nil {
		Logger.Error("Cannot send response", zap.Error(err))
		return err
	}

	return nil
}
