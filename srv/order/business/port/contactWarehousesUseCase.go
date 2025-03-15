package port

import (
	"context"
)

type IContactWarehousesUseCase interface {
	ContactWarehouses(context.Context, ContactWarehousesCmd) error
}

type ContactWarehousesCmd struct {
	OrderId               string
	LastContact           int64
	ConfirmedReservations []ConfirmedReservation
	ExcludeWarehouses     []string
}

type ConfirmedReservation struct {
	WarehouseId   string
	ReservationID string
	Goods         map[string]int64
}
