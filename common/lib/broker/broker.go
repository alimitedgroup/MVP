package broker

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
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
	*zap.Logger
}

func NewNatsConn(cfg *BrokerConfig, logger *zap.Logger) (*nats.Conn, error) {
	logger.Debug("Connecting to NATS", zap.String("url", cfg.Url))

	nc, err := nats.Connect(cfg.Url)
	if err != nil {
		return nil, err
	}

	return nc, nil
}

func NewNatsMessageBroker(nc *nats.Conn, logger *zap.Logger) (*NatsMessageBroker, error) {
	ncJs, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	return &NatsMessageBroker{Nats: nc, NatsJs: ncJs, Js: js, Logger: logger}, nil
}

func (n *NatsMessageBroker) RegisterRequest(ctx context.Context, subject Subject, queue Queue, handler RequestHandler) error {
	var sub *nats.Subscription

	sub, err := n.Nats.QueueSubscribe(subject.String(), queue.String(), func(msg *nats.Msg) {
		err := handler(ctx, msg)
		if err != nil {
			if errUnsub := sub.Unsubscribe(); errUnsub != nil {
				n.Fatal(
					"Error unsubscribing after another error",
					zap.Error(errUnsub),
					zap.NamedError("original_error", errUnsub),
					zap.String("subject", subject.String()),
				)
			}
			n.Fatal(
				"Error handling request",
				zap.Error(err),
				zap.String("subject", subject.String()),
				zap.String("queue", queue.String()),
			)
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
			n.Fatal(
				"Failed to handle message",
				zap.Error(msgErr),
				zap.String("subject", m.Subject()),
				zap.String("stream", streamCfg.Name),
			)
		} else {
			if errAck := m.Ack(); errAck != nil {
				cc.Stop()
				n.Fatal("Failed to ack message", zap.Error(errAck))
			}
		}

		var meta *jetstream.MsgMetadata
		meta, msgErr = m.Metadata()
		if msgErr != nil {
			cc.Stop()
			n.Fatal("Failed to read message metadata", zap.Error(msgErr))
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
