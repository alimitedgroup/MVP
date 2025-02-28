package broker

import (
	"context"
	"fmt"
	"log"

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

func (n *NatsMessageBroker) RegisterRequest(ctx context.Context, subject Subject, queue Queue, handler RequestHandler) error {
	var sub *nats.Subscription
	var err error

	sub, err = n.Nats.QueueSubscribe(subject.String(), queue.String(), func(msg *nats.Msg) {
		err = handler(ctx, msg)
		if err != nil {
			err = sub.Unsubscribe()
		}
	})
	if err != nil {
		return err
	}

	_ = sub

	return nil
}

func (n *NatsMessageBroker) RegisterJsHandler(ctx context.Context, restore *RestoreStreamControl, streamCfg jetstream.StreamConfig, handler JsHandler, opts ...JsHandlerOpt) error {
	s, err := n.Js.CreateStream(ctx, streamCfg)
	if err != nil {
		return err
	}

	consumerCfg := jetstream.ConsumerConfig{}
	for _, opt := range opts {
		opt(&consumerCfg)
	}

	consumer, err := s.CreateConsumer(ctx, consumerCfg)
	if err != nil {
		return err
	}
	restore.Start()

	// Consume all existing messages, and when they are finished unlock the waitgroup and continue listening
	var cc jetstream.ConsumeContext
	var msgErr error
	var isWgUnlocked bool = false

	cc, err = consumer.Consume(func(m jetstream.Msg) {
		msgErr = handler(ctx, m)
		if msgErr != nil {
			err = fmt.Errorf("failed to handle message: %w", msgErr)
			cc.Stop()
			return
		} else {
			err = m.Ack()
			if err != nil {
				err = fmt.Errorf("failed to ack message: %w", err)
				cc.Stop()
				return
			}
		}

		var meta *jetstream.MsgMetadata
		meta, msgErr = m.Metadata()
		if msgErr != nil {
			err = fmt.Errorf("failed to read message metadata: %w", msgErr)
			cc.Stop()
			return
		}

		if msgErr == nil && meta.NumPending == 0 && !isWgUnlocked {
			restore.Finish()
			isWgUnlocked = true
		}
	})

	return nil
}
