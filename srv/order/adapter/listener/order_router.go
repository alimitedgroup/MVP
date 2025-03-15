package listener

import (
	"context"
	"fmt"
	"log"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	internalStream "github.com/alimitedgroup/MVP/srv/order/adapter/stream"
	"github.com/nats-io/nats.go/jetstream"
)

type OrderRouter struct {
	orderListener *OrderListener
	broker        *broker.NatsMessageBroker
	restore       broker.IRestoreStreamControl
}

func NewOrderRouter(restoreFactory broker.IRestoreStreamControlFactory, orderListener *OrderListener, broker *broker.NatsMessageBroker) *OrderRouter {
	return &OrderRouter{orderListener, broker, restoreFactory.Build()}
}

func (r *OrderRouter) Setup(ctx context.Context) error {
	// register stream message handlers
	if err := r.broker.RegisterJsHandler(ctx, r.restore, stream.OrderUpdateStreamConfig, r.orderListener.ListenOrderUpdate); err != nil {
		return err
	}

	// wait restoring of the state before starting the server
	r.restore.Wait()

	if err := r.registerJsWithConsumerGroup(ctx,
		internalStream.ContactWarehousesStreamConfig,
		internalStream.ContactWarehousesStreamConsumerConfig,
		r.orderListener.ListenContactWarehouses,
	); err != nil {
		return err
	}

	// register request/reply handlers

	return nil
}

func (r *OrderRouter) registerJsWithConsumerGroup(ctx context.Context, streamCfg jetstream.StreamConfig, consumerCfg jetstream.ConsumerConfig, handler broker.JsHandler) error {
	s, err := r.broker.Js.CreateStream(ctx, streamCfg)
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
