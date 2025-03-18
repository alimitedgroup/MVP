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
	Nats   *nats.Conn
	NatsJs nats.JetStream
	Js     jetstream.JetStream
}

func NewNatsConn(cfg *BrokerConfig) (*nats.Conn, error) {
	log.Printf("Connecting to NATS at %s\n", cfg.Url)

	nc, err := nats.Connect(cfg.Url)
	if err != nil {
		return nil, err
	}

	return nc, nil
}

func NewNatsMessageBroker(nc *nats.Conn) (*NatsMessageBroker, error) {
	ncJs, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	return &NatsMessageBroker{nc, ncJs, js}, nil
}

func (n *NatsMessageBroker) RegisterRequest(ctx context.Context, subject Subject, queue Queue, handler RequestHandler) error {
	var sub *nats.Subscription

	sub, err := n.Nats.QueueSubscribe(subject.String(), queue.String(), func(msg *nats.Msg) {
		err := handler(ctx, msg)
		if err != nil {
			if errUnsub := sub.Unsubscribe(); errUnsub != nil {
				log.Fatalf("Error unsubscribing: %v\nafter error %v\n", errUnsub, err)
			}
			log.Fatalf("Error handling request: %v\n", err)
		}
	})
	if err != nil {
		return err
	}

	_ = sub

	return nil
}

func (n *NatsMessageBroker) RegisterJsHandler(ctx context.Context, restore IRestoreStreamControl, streamCfg jetstream.StreamConfig, handler JsHandler, opts ...JsHandlerOpt) error {
	s, err := n.Js.CreateStream(ctx, streamCfg)
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}

	consumerCfg := jetstream.ConsumerConfig{}
	for _, opt := range opts {
		opt(&consumerCfg)
	}

	consumer, err := s.CreateConsumer(ctx, consumerCfg)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	var isWgUnlocked bool = false
	restore.Start()

	// fetch consumer info
	info, err := consumer.Info(ctx)
	if err != nil {
		return fmt.Errorf("failed to get consumer info: %w", err)
	}
	// if num pending is zero the stream have just been created and there are no messages to consume
	if info.NumPending == 0 {
		restore.Finish()
		isWgUnlocked = true
	}

	// Consume all existing messages, and when they are finished unlock the waitgroup and continue listening
	var cc jetstream.ConsumeContext

	cc, err = consumer.Consume(func(m jetstream.Msg) {
		msgErr := handler(ctx, m)
		if msgErr != nil {
			cc.Stop()
			log.Fatalf("failed to handle message: %v\n", msgErr)
		} else {
			if errAck := m.Ack(); errAck != nil {
				cc.Stop()
				log.Fatalf("failed to ack message: %v\nafter error: %v\n", errAck, err)
			}
		}

		var meta *jetstream.MsgMetadata
		meta, msgErr = m.Metadata()
		if msgErr != nil {
			cc.Stop()
			log.Fatalf("failed to read message metadata: %v\n", msgErr)
		}

		if meta.NumPending == 0 && !isWgUnlocked {
			restore.Finish()
			isWgUnlocked = true
		}
	})

	if err != nil {
		return fmt.Errorf("failed to consume messages: %w", err)
	}

	return nil
}

func (n *NatsMessageBroker) RegisterJsWithConsumerGroup(ctx context.Context, streamCfg jetstream.StreamConfig, consumerCfg jetstream.ConsumerConfig, handler JsHandler) error {
	s, err := n.Js.CreateStream(ctx, streamCfg)
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}

	consumer, err := s.CreateOrUpdateConsumer(ctx, consumerCfg)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	var cc jetstream.ConsumeContext

	cc, err = consumer.Consume(func(m jetstream.Msg) {
		msgErr := handler(ctx, m)
		if msgErr != nil {
			if msgErr == ErrMsgNotAcked {
				return
			}
			cc.Stop()
			log.Fatalf("failed to handle message: %v\n", msgErr)
		} else {
			if errAck := m.Ack(); errAck != nil {
				cc.Stop()
				log.Fatalf("failed to ack message: %v\nafter error: %v\n", errAck, err)
			}
		}
	})
	if err != nil {
		return fmt.Errorf("failed to consume messages: %w", err)
	}

	return nil

}

var ErrMsgNotAcked = fmt.Errorf("message not acked")
