package controller

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
)

type TransferController struct {
	broker                *broker.NatsMessageBroker
	createTransferUseCase port.ICreateTransferUseCase
	getTransferUseCase    port.IGetTransferUseCase
}

type TransferControllerParams struct {
	fx.In

	CreateTransferUseCase port.ICreateTransferUseCase
	GetTransferUseCase    port.IGetTransferUseCase
	Broker                *broker.NatsMessageBroker
}

func NewTransferController(p TransferControllerParams) *TransferController {
	return &TransferController{p.Broker, p.CreateTransferUseCase, p.GetTransferUseCase}
}

func (c *TransferController) TransferCreateHandler(ctx context.Context, msg *nats.Msg) error {
	var dto request.CreateTransferRequestDTO
	if err := json.Unmarshal(msg.Data, &dto); err != nil {
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
		SenderId:   dto.SenderId,
		ReceiverId: dto.ReceiverId,
		Goods:      goods,
	}
	resp, err := c.createTransferUseCase.CreateTransfer(ctx, cmd)
	if err != nil {
		if err := broker.RespondToMsg(msg, response.ErrorResponseDTO{Error: err.Error()}); err != nil {
			return err
		}
	}

	respDto := response.TransferCreateResponseDTO{
		Message: response.TransferCreateInfo{TransferID: resp.TransferID},
	}
	if err := broker.RespondToMsg(msg, respDto); err != nil {
		return err
	}

	return nil
}

func (c *TransferController) TransferGetHandler(ctx context.Context, msg *nats.Msg) error {
	var dto request.GetTransferRequestDTO
	if err := json.Unmarshal(msg.Data, &dto); err != nil {
		return err
	}

	transfer, err := c.getTransferUseCase.GetTransfer(ctx, dto.TransferID)
	if err != nil {
		return err
	}

	respDto := response.GetTransferResponseDTO{
		Message: modelTransferToTransferInfoDTO(transfer),
	}
	if err := broker.RespondToMsg(msg, respDto); err != nil {
		return err
	}

	return nil
}

func (c *TransferController) TransferGetAllHandler(ctx context.Context, msg *nats.Msg) error {
	transfers, err := c.getTransferUseCase.GetAllTransfers(ctx)
	if err != nil {
		return err
	}

	transfersDto := make([]response.TransferInfo, 0, len(transfers))
	for _, transfer := range transfers {
		transfersDto = append(transfersDto, modelTransferToTransferInfoDTO(transfer))
	}
	respDto := response.GetAllTransferResponseDTO{
		Message: transfersDto,
	}
	if err := broker.RespondToMsg(msg, respDto); err != nil {
		return err
	}

	return nil
}

func modelTransferToTransferInfoDTO(transfer model.Transfer) response.TransferInfo {
	goods := make([]response.TransferInfoGood, 0, len(transfer.Goods))
	for _, good := range transfer.Goods {
		goods = append(goods, response.TransferInfoGood{
			GoodID:   string(good.ID),
			Quantity: good.Quantity,
		})
	}
	return response.TransferInfo{
		Status:     transfer.Status,
		TransferID: string(transfer.Id),
		SenderID:   string(transfer.SenderId),
		ReceiverID: string(transfer.ReceiverId),
		Goods:      goods,
	}
}
