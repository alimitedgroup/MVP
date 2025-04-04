package stream

import "github.com/nats-io/nats.go/jetstream"

var ReservationEventStreamConfig = jetstream.StreamConfig{
	Name:     "reservation",
	Subjects: []string{"reservation.>"},
}

type ReservationEvent struct {
	ID          string            `json:"id"`
	WarehouseID string            `json:"warehouse_id"`
	Goods       []ReservationGood `json:"goods"`
}

type ReservationGood struct {
	GoodID   string `json:"good_id"`
	Quantity int64  `json:"quantity"`
}
