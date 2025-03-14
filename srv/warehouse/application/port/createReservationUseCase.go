package port

import "context"

type ICreateReservationUseCase interface {
	CreateReservation(context.Context, CreateReservationCmd) error
}

type CreateReservationCmd struct {
}
