package controller

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	GetAllTransferRequestCounter metric.Int64Counter
	GetTransferRequestCounter    metric.Int64Counter
	TransferCreateRequestCounter metric.Int64Counter
)

type TransferController struct {
	createTransferUseCase port.ICreateTransferUseCase
	getTransferUseCase    port.IGetTransferUseCase
}

type TransferControllerParams struct {
	fx.In

	CreateTransferUseCase port.ICreateTransferUseCase
	GetTransferUseCase    port.IGetTransferUseCase
	Logger                *zap.Logger
	Meter                 metric.Meter
}

func NewTransferController(p TransferControllerParams) *TransferController {
	observability.CounterSetup(&p.Meter, p.Logger, &TotalRequestCounter, &MetricMap, "num_order_transfer_requests")
	observability.CounterSetup(&p.Meter, p.Logger, &TransferCreateRequestCounter, &MetricMap, "num_transfer_create_requests")
	observability.CounterSetup(&p.Meter, p.Logger, &GetTransferRequestCounter, &MetricMap, "num_get_transfer_requests")
	observability.CounterSetup(&p.Meter, p.Logger, &GetAllTransferRequestCounter, &MetricMap, "num_get_all_transfer_requests")
	Logger = p.Logger
	return &TransferController{p.CreateTransferUseCase, p.GetTransferUseCase}
}

func (c *TransferController) TransferCreateHandler(ctx context.Context, msg *nats.Msg) error {

	Logger.Info("Received new transfer creation request")
	verdict := "success"

	defer func() {
		Logger.Info("Completed create transfer request")
		TotalRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		TransferCreateRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	var dto request.CreateTransferRequestDTO
	if err := json.Unmarshal(msg.Data, &dto); err != nil {
		Logger.Debug("Bad request", zap.Error(err))
		verdict = "bad request"
		return err
	}

	goods := make([]port.CreateTransferGood, 0, len(dto.Goods))
	for _, good := range dto.Goods {
		goods = append(goods, port.CreateTransferGood{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}
	cmd := port.CreateTransferCmd{
		SenderID:   dto.SenderID,
		ReceiverID: dto.ReceiverID,
		Goods:      goods,
	}
	resp, err := c.createTransferUseCase.CreateTransfer(ctx, cmd)
	if err != nil {
		verdict = "cannot create order"
		Logger.Debug("Cannot create order", zap.Error(err))
		if err := broker.RespondToMsg(msg, response.ErrorResponseDTO{Error: err.Error()}); err != nil {
			Logger.Error("Cannot send response", zap.Error(err))
			return err
		}
	}

	respDto := response.TransferCreateResponseDTO{
		Message: response.TransferCreateInfo{TransferID: resp.TransferID},
	}
	if err := broker.RespondToMsg(msg, respDto); err != nil {
		Logger.Error("Cannot send response", zap.Error(err))
		return err
	}

	return nil
}

func (c *TransferController) TransferGetHandler(ctx context.Context, msg *nats.Msg) error {

	Logger.Info("Received new get transfer request")
	verdict := "success"

	defer func() {
		Logger.Info("Completed get transfer request")
		TotalRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		GetTransferRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	var dto request.GetTransferRequestDTO
	if err := json.Unmarshal(msg.Data, &dto); err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		return err
	}

	transfer, err := c.getTransferUseCase.GetTransfer(ctx, port.GetTransferCmd(dto.TransferID))
	if err != nil {
		verdict = "cannot get transfer"
		Logger.Debug("Cannot get transfer", zap.Error(err))
		return err
	}

	respDto := response.GetTransferResponseDTO{
		Message: modelTransferToTransferInfoDTO(transfer),
	}
	if err := broker.RespondToMsg(msg, respDto); err != nil {
		Logger.Error("Cannot send response", zap.Error(err))
		return err
	}

	return nil
}

func (c *TransferController) TransferGetAllHandler(ctx context.Context, msg *nats.Msg) error {

	Logger.Info("Received new get all transfer request")
	verdict := "success"

	defer func() {
		Logger.Info("Completed get all transfer request")
		TotalRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		GetAllTransferRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	transfers := c.getTransferUseCase.GetAllTransfers(ctx)

	transfersDto := make([]response.TransferInfo, 0, len(transfers))
	for _, transfer := range transfers {
		transfersDto = append(transfersDto, modelTransferToTransferInfoDTO(transfer))
	}
	respDto := response.GetAllTransferResponseDTO{
		Message: transfersDto,
	}
	if err := broker.RespondToMsg(msg, respDto); err != nil {
		Logger.Error("Cannot send response", zap.Error(err))
		return err
	}

	return nil
}

func modelTransferToTransferInfoDTO(transfer model.Transfer) response.TransferInfo {
	goods := make([]response.TransferInfoGood, 0, len(transfer.Goods))
	for _, good := range transfer.Goods {
		goods = append(goods, response.TransferInfoGood{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}
	return response.TransferInfo{
		Status:       transfer.Status,
		TransferID:   transfer.ID,
		SenderID:     transfer.SenderID,
		ReceiverID:   transfer.ReceiverID,
		CreationTime: transfer.CreationTime,
		UpdateTime:   transfer.UpdateTime,
		Goods:        goods,
	}
}
