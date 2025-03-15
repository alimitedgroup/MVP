package listener

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/stream"
)

type ReservationEventRouter struct {
	reservationListener *ReservationEventListener
	broker              *broker.NatsMessageBroker
	restore             broker.IRestoreStreamControl
}

func NewReservationEventRouter(restoreFactory broker.IRestoreStreamControlFactory, reservationListener *ReservationEventListener, n *broker.NatsMessageBroker) *ReservationEventRouter {
	return &ReservationEventRouter{reservationListener, n, restoreFactory.Build()}
}

func (r *ReservationEventRouter) Setup(ctx context.Context) error {
	// register stream message handlers
	err := r.broker.RegisterJsHandler(ctx, r.restore, stream.ReservationEventStreamConfig, r.reservationListener.ListenReservationEvent)
	if err != nil {
		return err
	}

	// wait restoring of the state before starting the server
	r.restore.Wait()

	// register request/reply handlers

	return nil
}
