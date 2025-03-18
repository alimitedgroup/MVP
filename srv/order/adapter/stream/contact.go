package stream

import (
	"github.com/nats-io/nats.go/jetstream"
)

var ContactWarehousesStreamConfig = jetstream.StreamConfig{
	Name:     "order_contact_warehouses",
	Subjects: []string{"order.contact.warehouses"},
	Storage:  jetstream.FileStorage,
}

type ContactWarehouses struct {
	Order                 *ContactWarehousesOrder
	Transfer              *ContactWarehousesTransfer
	LastContact           int64
	ConfirmedReservations []ConfirmedReservation
	ExcludeWarehouses     []string
	Type                  ContactWarehousesType
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
	SenderId      string
	ReceiverId    string
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

var ContactWarehousesStreamConsumerConfig = jetstream.ConsumerConfig{
	Durable: "order_contact_warehouses",
}
