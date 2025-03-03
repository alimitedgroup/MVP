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

const ApiGatewayQueue Queue = "api_gateway"
const NoQueue Queue = ""

var StockUpdateSubject Subject = "stock.update"

type JsHandler func(context.Context, jetstream.Msg) error

// JsHandlerOpt represents various options used when creating a JetStream handler
type JsHandlerOpt func(config *jetstream.ConsumerConfig)

// WithDeliverNew will set the consumer's DeliveryPolicy to DeliverNew
func WithDeliverNew() JsHandlerOpt {
	return func(config *jetstream.ConsumerConfig) {
		config.DeliverPolicy = jetstream.DeliverNewPolicy
	}
}

// WithDeliverAll will set the consumer's DeliveryPolicy to DeliverAll
func WithDeliverAll() JsHandlerOpt {
	return func(config *jetstream.ConsumerConfig) {
		config.DeliverPolicy = jetstream.DeliverAllPolicy
	}
}

// WithSubjectFilter will filter the delivered messages to those specified. Mutually exclusive with WithSubjectsFilter
func WithSubjectFilter(subject string) JsHandlerOpt {
	return func(config *jetstream.ConsumerConfig) {
		config.FilterSubject = subject
	}
}

// WithSubjectsFilter will filter the delivered messages to those specified. Mutually exclusive with WithSubjectFilter
func WithSubjectsFilter(subjects []string) JsHandlerOpt {
	return func(config *jetstream.ConsumerConfig) {
		config.FilterSubjects = append(config.FilterSubjects, subjects...)
	}
}
