package port

import "context"

type ICreateReservationEventUseCase interface {
	CreateReservationEvent(context.Context, CreateReservationEventCmd) error
}

type CreateReservationEventCmd struct {
}
