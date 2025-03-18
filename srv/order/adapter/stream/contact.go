package stream

import (
	"github.com/nats-io/nats.go/jetstream"
)

var ContactWarehousesStreamConfig = jetstream.StreamConfig{
	Name:     "contact_warehouses",
	Subjects: []string{"contact.warehouses"},
	Storage:  jetstream.FileStorage,
}

type ContactWarehouses struct {
	Order                 *ContactWarehousesOrder    `json:"order,omitempty"`
	Transfer              *ContactWarehousesTransfer `json:"transfer,omitempty"`
	ConfirmedReservations []ConfirmedReservation     `json:"confirmed_reservations"`
	ExcludeWarehouses     []string                   `json:"exclude_warehouses"`
	Type                  ContactWarehousesType      `json:"type"`
	RetryInTime           int64                      `json:"retry_in_time"`
	RetryUntil            int64                      `json:"retry_until"`
}

type ContactWarehousesType string

var (
	ContactWarehousesTypeOrder    ContactWarehousesType = "order"
	ContactWarehousesTypeTransfer ContactWarehousesType = "transfer"
)

type ContactWarehousesOrder struct {
	ID           string                  `json:"id"`
	Status       string                  `json:"status"`
	Name         string                  `json:"name"`
	FullName     string                  `json:"full_name"`
	Address      string                  `json:"address"`
	UpdateTime   int64                   `json:"update_time"`
	CreationTime int64                   `json:"creation_time"`
	Goods        []ContactWarehousesGood `json:"goods"`
	Reservations []string                `json:"reservations"`
}

type ContactWarehousesTransfer struct {
	ID            string                  `json:"id"`
	Status        string                  `json:"status"`
	SenderId      string                  `json:"sender_id"`
	ReceiverId    string                  `json:"receiver_id"`
	UpdateTime    int64                   `json:"update_time"`
	CreationTime  int64                   `json:"creation_time"`
	Goods         []ContactWarehousesGood `json:"goods"`
	ReservationId string                  `json:"reservation_id"`
}

type ContactWarehousesGood struct {
	GoodID   string `json:"good_id"`
	Quantity int64  `json:"quantity"`
}

type ConfirmedReservation struct {
	WarehouseId   string           `json:"warehouse_id"`
	ReservationID string           `json:"reservation_id"`
	Goods         map[string]int64 `json:"goods"`
}

var ContactWarehousesStreamConsumerConfig = jetstream.ConsumerConfig{
	Durable:       "contact_warehouses",
	AckPolicy:     jetstream.AckExplicitPolicy,
	DeliverPolicy: jetstream.DeliverAllPolicy,
	// MaxDeliver:    100,
	// AckWait:       time.Duration(24) * time.Hour,
	// BackOff:       []time.Duration{10 * time.Second, 20 * time.Second, 40 * time.Second},
}
