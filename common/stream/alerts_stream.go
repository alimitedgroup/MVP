package stream

import "github.com/nats-io/nats.go/jetstream"

var AlertConfig = jetstream.StreamConfig{
	Name:              "alerts",
	Subjects:          []string{"stock.alert.>"},
	MaxMsgsPerSubject: 1,
}
