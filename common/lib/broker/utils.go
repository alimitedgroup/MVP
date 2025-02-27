package broker

import (
	"github.com/nats-io/nats.go"
)

type Queue string

func (q Queue) String() string {
	return string(q)
}

type Subject string

func (s Subject) String() string {
	return string(s)
}

type RequestHandler func(msg *nats.Msg)

const ApiGatewayQueue Queue = "api_gateway"
const NoQueue Queue = ""

var StockUpdateSubject Subject = "stock.update"
