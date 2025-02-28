package broker

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Queue string

func (q Queue) String() string {
	return string(q)
}

type Subject string

func (s Subject) String() string {
	return string(s)
}

type RequestHandler func(context.Context, *nats.Msg) error

type JsHandler func(context.Context, jetstream.Msg) error

const ApiGatewayQueue Queue = "api_gateway"
const NoQueue Queue = ""

var StockUpdateSubject Subject = "stock.update"
