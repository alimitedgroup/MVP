package controller

import (
	"context"
	"fmt"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
)

type ReservationRouter struct {
	config                *config.WarehouseConfig
	reservationController *ReservationController
	broker                *broker.NatsMessageBroker
}

func NewReservationRouter(config *config.WarehouseConfig, reservationController *ReservationController, broker *broker.NatsMessageBroker) *ReservationRouter {
	return &ReservationRouter{config, reservationController, broker}
}

func (r *ReservationRouter) Setup(ctx context.Context) error {
	// register request/reply handlers
	err := r.broker.RegisterRequest(ctx, broker.Subject(fmt.Sprintf("warehouse.reservation.add.%s", r.config.ID)), broker.NoQueue, r.reservationController.CreateReservationHandler)
	if err != nil {
		return err
	}

	return nil
}
