package port

import (
	"context"
)

type IRequestReservationPort interface {
	RequestReservation(context.Context, RequestReservationCmd) (RequestReservationResponse, error)
}

type RequestReservationCmd struct {
	WarehouseId string
	Items       []ReservationItem
}

type ReservationItem struct {
	GoodID   string
	Quantity int64
}

type RequestReservationResponse struct {
	Id string
}
