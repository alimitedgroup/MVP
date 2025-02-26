package broker

import (
	"strings"

	"github.com/nats-io/nats.go"
)

type Queue string

func (q Queue) String() string {
	return string(q)
}

const ApiGatewayQueue Queue = "api_gateway"
const NoQueue Queue = ""

type Subject []string

func NewSubject() Subject {
	return Subject{}
}

var StockUpdateSubject Subject = NewSubject().S("stock").S("update")

func (s Subject) S(name string) Subject {
	s = append(s, name)
	return s
}

func (s Subject) All() Subject {
	s = append(s, "*")
	return s
}

func (s Subject) Name() string {
	return strings.Join(s, ".")
}

type RequestHandler func(msg *nats.Msg)
