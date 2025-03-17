package port

import (
	"context"
)

type IContactWarehousesUseCase interface {
	ContactWarehouses(context.Context, ContactWarehousesCmd) error
}

type ContactWarehousesCmd struct {
	Order                 ContactWarehousesOrder
	LastContact           int64
	ConfirmedReservations []ConfirmedReservation
	ExcludeWarehouses     []string
}

type ContactWarehousesOrder struct {
	ID           string
	Status       string
	Name         string
	Email        string
	Address      string
	UpdateTime   int64
	CreationTime int64
	Goods        []ContactWarehousesGood
	Reservations []string
}

type ContactWarehousesGood struct {
	GoodId   string
	Quantity int64
}

type ConfirmedReservation struct {
	WarehouseId   string
	ReservationID string
	Goods         map[string]int64
}
