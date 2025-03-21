package port

import (
	"context"
	"errors"
)

type ICreateReservationUseCase interface {
	CreateReservation(context.Context, CreateReservationCmd) (CreateReservationResponse, error)
}

type CreateReservationCmd struct {
	Goods []ReservationGood
}
type ReservationGood struct {
	GoodID   string
	Quantity int64
}

type CreateReservationResponse struct {
	ReservationID string
}

var ErrNotEnoughStock = errors.New("not enough stock")
