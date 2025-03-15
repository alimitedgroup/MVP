package stream

import (
	"github.com/nats-io/nats.go/jetstream"
)

var OrderUpdateStreamConfig = jetstream.StreamConfig{
	Name:     "order_update",
	Subjects: []string{"order.update"},
	Storage:  jetstream.FileStorage,
}

type OrderUpdate struct {
	ID           string            `json:"id"`
	Status       string            `json:"status"`
	Name         string            `json:"name"`
	Email        string            `json:"email"`
	Address      string            `json:"address"`
	Goods        []OrderUpdateGood `json:"goods"`
	CreationTime int64             `json:"creation_time"`
	UpdateTime   int64             `json:"update_time"`
}

type OrderUpdateGood struct {
	GoodID   string `json:"id"`
	Quantity int64  `json:"quantity"`
}
