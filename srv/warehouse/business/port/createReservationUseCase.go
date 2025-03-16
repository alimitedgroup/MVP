package port

import "context"

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
