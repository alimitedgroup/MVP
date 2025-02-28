package broker

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type MessageBroker interface {
}

type BrokerConfig struct {
	Url string `mapstructure:"url"`
}

type NatsMessageBroker struct {
	Nats *nats.Conn
	Js   jetstream.JetStream
}

func NewNatsMessageBroker(cfg *BrokerConfig) (*NatsMessageBroker, error) {
	log.Printf("Connecting to NATS at %s\n", cfg.Url)

	nc, err := nats.Connect(cfg.Url)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	return &NatsMessageBroker{nc, js}, nil
}

func (n *NatsMessageBroker) RequestSubscribe(ctx context.Context, subject Subject, queue Queue, handler RequestHandler) (*nats.Subscription, error) {
	sub, err := n.Nats.QueueSubscribe(subject.String(), queue.String(), func(msg *nats.Msg) {
		handler(ctx, msg)
	})
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (n *NatsMessageBroker) RegisterJsHandler(ctx context.Context, subject Subject, handler JsHandler, wg *sync.WaitGroup) error {
	n.Js.CreateStream(ctx, jetstream.StreamConfig{
		Name: subject.String(),
	})

	cfg := jetstream.ConsumerConfig{}
	consumer, err := n.Js.CreateConsumer(ctx, subject.String(), cfg)
	if err != nil {
		return err
	}
	wg.Add(1)

	// Consume all existing messages, and when they are finished unlock the waitgroup and continue listening
	var cc jetstream.ConsumeContext
	var msgErr error
	var isWgUnlocked bool = false

	cc, err = consumer.Consume(func(m jetstream.Msg) {
		msgErr = handler(ctx, m)
		if msgErr != nil {
			err = fmt.Errorf("failed to handle message: %w", msgErr)
			cc.Stop()
		}

		var meta *jetstream.MsgMetadata
		meta, msgErr = m.Metadata()
		if msgErr != nil {
			err = fmt.Errorf("failed to read message metadata: %w", msgErr)
			cc.Stop()
		}
		if msgErr == nil && meta.NumPending == 0 && !isWgUnlocked {
			wg.Done()
			isWgUnlocked = true
		}
	})

	return nil
}
