package constants

import (
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/nats-io/nats.go/jetstream"
)

var ApiGatewayGroup broker.Queue = "api_gateway"

var StockUpdatesStreamConfig = jetstream.StreamConfig{
	Name:     "stock_updates",
	Subjects: []string{"stock.updates.>"},
	Storage:  jetstream.FileStorage,
}
