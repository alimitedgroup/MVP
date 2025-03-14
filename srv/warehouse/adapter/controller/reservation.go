package controller

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/nats-io/nats.go"
)

type ReservationController struct {
}

func NewReservationController() *ReservationController {
	return &ReservationController{}
}

func (c *ReservationController) CreateReservationHandler(ctx context.Context, msg *nats.Msg) error {
	var dto request.ReserveStockRequestDTO
	if err := json.Unmarshal(msg.Data, &dto); err != nil {
		return err
	}

	reservationId := "reservation_id"

	respDto := response.ReserveStockResponseDTO{
		Message: response.ReserveStockInfo{
			ReservationID: reservationId,
		},
	}
	if err := broker.RespondToMsg(msg, &respDto); err != nil {
		return err
	}

	return nil
}
