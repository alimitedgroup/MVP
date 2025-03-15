package stream

import "github.com/nats-io/nats.go/jetstream"

var ContactWarehousesStreamConfig = jetstream.StreamConfig{
	Name:     "order_contact_warehouses",
	Subjects: []string{"order.contact.warehouses"},
	Storage:  jetstream.FileStorage,
}

type ContactWarehouses struct {
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

var ContactWarehousesStreamConsumerConfig = jetstream.ConsumerConfig{
	Durable: "order_contact_warehouses",
}
