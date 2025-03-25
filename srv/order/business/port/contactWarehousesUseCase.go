package port

import (
	"context"
	"time"
)

type IContactWarehousesUseCase interface {
	ContactWarehouses(context.Context, ContactWarehousesCmd) (ContactWarehousesResponse, error)
}

type ContactWarehousesCmd struct {
	Type                  ContactWarehousesType
	Order                 *ContactWarehousesOrder
	Transfer              *ContactWarehousesTransfer
	ConfirmedReservations []ConfirmedReservation
	ExcludeWarehouses     []string
	RetryUntil            int64
	RetryInTime           int64
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
	ReservationID string
}

type ContactWarehousesGood struct {
	GoodID   string
	Quantity int64
}

type ConfirmedReservation struct {
	WarehouseID   string
	ReservationID string
	Goods         map[string]int64
}

type ContactWarehousesResponse struct {
	IsRetry    bool
	RetryAfter time.Duration
}
