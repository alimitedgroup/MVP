package stream

import (
	"github.com/nats-io/nats.go/jetstream"
)

var StockUpdateStreamConfig = jetstream.StreamConfig{
	Name:     "stock_update",
	Subjects: []string{"stock.update.>"},
	Storage:  jetstream.FileStorage,
}

type StockUpdate struct {
	ID            string            `json:"id"`
	WarehouseID   string            `json:"warehouse_id"`
	Type          StockUpdateType   `json:"type"`
	Goods         []StockUpdateGood `json:"goods"`
	OrderID       string            `json:"order_id"`
	TransferID    string            `json:"transfer_id"`
	ReservationID string            `json:"reservation_id"`
	Timestamp     int64             `json:"timestamp"`
}

type StockUpdateType string

const (
	StockUpdateTypeAdd    StockUpdateType = "add"
	StockUpdateTypeRemove StockUpdateType = "remove"
	StockUpdateTypeOrder  StockUpdateType = "order"
)

type StockUpdateGood struct {
	GoodID   string `json:"id"`
	Quantity int64  `json:"quantity"`
	Delta    int64  `json:"delta"`
}
