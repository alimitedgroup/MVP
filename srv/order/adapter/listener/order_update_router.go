package listener

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
)

type OrderUpdateRouter struct {
	orderListener *OrderUpdateListener
	broker        *broker.NatsMessageBroker
	restore       broker.IRestoreStreamControl
}

func NewOrderUpdateRouter(restoreFactory broker.IRestoreStreamControlFactory, orderUpdateListener *OrderUpdateListener, broker *broker.NatsMessageBroker) *OrderUpdateRouter {
	return &OrderUpdateRouter{orderUpdateListener, broker, restoreFactory.Build()}
}

func (r *OrderUpdateRouter) Setup(ctx context.Context) error {
	// register stream message handlers
	err := r.broker.RegisterJsHandler(ctx, r.restore, stream.OrderUpdateStreamConfig, r.orderListener.ListenOrderUpdate)
	if err != nil {
		return err
	}

	// wait restoring of the state before starting the server
	r.restore.Wait()

	// register request/reply handlers

	return nil
}
