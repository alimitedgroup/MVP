package controller

import (
	"context"

	"github.com/nats-io/nats.go"
)

type ReservationController struct {
}

func NewReservationController() *ReservationController {
	return &ReservationController{}
}

func (c *ReservationController) CreateReservationHandler(ctx context.Context, msg *nats.Msg) error {
	return nil
}
