package listener

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
)

type OrderUpdateRouter struct {
	reservationListener *OrderUpdateListener
	broker              *broker.NatsMessageBroker
	restore             broker.IRestoreStreamControl
}

func NewOrderUpdateRouter(restoreFactory broker.IRestoreStreamControlFactory, reservationListener *OrderUpdateListener, n *broker.NatsMessageBroker) *OrderUpdateRouter {
	return &OrderUpdateRouter{reservationListener, n, restoreFactory.Build()}
}

func (r *OrderUpdateRouter) Setup(ctx context.Context) error {
	// register stream message handlers
	err := r.broker.RegisterJsHandler(ctx, r.restore, stream.OrderUpdateStreamConfig, r.reservationListener.ListenOrderUpdate)
	if err != nil {
		return err
	}

	// wait restoring of the state before starting the server
	r.restore.Wait()

	// register request/reply handlers

	return nil
}
