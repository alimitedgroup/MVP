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

type ReservationController struct {
	createReservationUseCase port.ICreateReservationUseCase
}

func NewReservationController(createReservationUseCase port.ICreateReservationUseCase) *ReservationController {
	return &ReservationController{createReservationUseCase}
}

func (c *ReservationController) CreateReservationHandler(ctx context.Context, msg *nats.Msg) error {
	var dto request.ReserveStockRequestDTO
	if err := json.Unmarshal(msg.Data, &dto); err != nil {
		return err
	}

	goods := make([]port.ReservationGood, 0, len(dto.Goods))
	for _, good := range dto.Goods {
		goods = append(goods, port.ReservationGood(good))
	}

	cmd := port.CreateReservationCmd{Goods: goods}
	createResp, err := c.createReservationUseCase.CreateReservation(ctx, cmd)
	if err != nil {
		resp := response.ErrorResponseDTO{
			Error: err.Error(),
		}
		if err := broker.RespondToMsg(msg, resp); err != nil {
			return err
		}
	}

	respDto := response.ReserveStockResponseDTO{
		Message: response.ReserveStockInfo{
			ReservationID: createResp.ReservationID,
		},
	}
	if err := broker.RespondToMsg(msg, &respDto); err != nil {
		return err
	}

	return nil
}
