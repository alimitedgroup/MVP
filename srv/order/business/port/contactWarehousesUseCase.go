package port

import (
	"context"
)

type IContactWarehousesUseCase interface {
	ContactWarehouses(context.Context, ContactWarehousesCmd) error
}

type ContactWarehousesCmd struct {
	Type                  ContactWarehousesType
	Order                 *ContactWarehousesOrder
	Transfer              *ContactWarehousesTransfer
	LastContact           int64
	ConfirmedReservations []ConfirmedReservation
	ExcludeWarehouses     []string
}

type ContactWarehousesType string

var (
	ContactWarehousesTypeOrder    ContactWarehousesType = "order"
	ContactWarehousesTypeTransfer ContactWarehousesType = "transfer"
)

type ContactWarehousesOrder struct {
	ID           string
	Status       string
	Name         string
	FullName     string
	Address      string
	UpdateTime   int64
	CreationTime int64
	Goods        []ContactWarehousesGood
	Reservations []string
}

type ContactWarehousesTransfer struct {
	ID            string
	Status        string
	SenderID      string
	ReceiverID    string
	UpdateTime    int64
	CreationTime  int64
	Goods         []ContactWarehousesGood
	ReservationId string
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
