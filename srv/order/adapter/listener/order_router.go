package listener

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	internalStream "github.com/alimitedgroup/MVP/srv/order/adapter/stream"
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

	if err := r.broker.RegisterJsWithConsumerGroup(ctx,
		internalStream.ContactWarehousesStreamConfig,
		internalStream.ContactWarehousesStreamConsumerConfig,
		r.orderListener.ListenContactWarehouses,
	); err != nil {
		return err
	}

	return nil
}
