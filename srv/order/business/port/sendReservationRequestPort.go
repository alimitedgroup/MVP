package port

import (
	"context"
)

type IRequestReservationPort interface {
	RequestReservation(context.Context, RequestReservationCmd) (RequestReservationResponse, error)
}

type RequestReservationCmd struct {
	WarehouseID string
	Goods       []ReservationGood
}

type ReservationGood struct {
	GoodID   string
	Quantity int64
}

type RequestReservationResponse struct {
	ID string
}
