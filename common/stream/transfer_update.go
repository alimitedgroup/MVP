package stream

import (
	"github.com/nats-io/nats.go/jetstream"
)

var TransferUpdateStreamConfig = jetstream.StreamConfig{
	Name:     "transfer_update",
	Subjects: []string{"transfer.update"},
	Storage:  jetstream.FileStorage,
}

type TransferUpdate struct {
	ID            string               `json:"id"`
	Status        string               `json:"status"`
	SenderID      string               `json:"sender_id"`
	ReceiverID    string               `json:"receiver_id"`
	Goods         []TransferUpdateGood `json:"goods"`
	ReservationID string               `json:"reservation_id"`
	CreationTime  int64                `json:"creation_time"`
	UpdateTime    int64                `json:"update_time"`
}

type TransferUpdateGood struct {
	GoodID   string `json:"id"`
	Quantity int64  `json:"quantity"`
}
