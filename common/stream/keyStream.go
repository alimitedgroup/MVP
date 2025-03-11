package stream

import "github.com/nats-io/nats.go/jetstream"

var KeyStream = jetstream.StreamConfig{
	Name:     "auth_keys",
	Subjects: []string{"keys.>"},
	Storage:  jetstream.FileStorage,
}
