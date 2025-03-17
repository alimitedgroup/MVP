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
	Order                 ContactWarehousesOrder
	TransferId            string
	LastContact           int64
	ConfirmedReservations []ConfirmedReservation
	ExcludeWarehouses     []string
}

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
