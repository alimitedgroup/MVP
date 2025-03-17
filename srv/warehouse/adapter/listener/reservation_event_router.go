package listener

import (
	"context"
	"fmt"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
)

type ReservationEventRouter struct {
	reservationListener *ReservationEventListener
	broker              *broker.NatsMessageBroker
	cfg                 *config.WarehouseConfig
	restore             broker.IRestoreStreamControl
}

func NewReservationEventRouter(restoreFactory broker.IRestoreStreamControlFactory, reservationListener *ReservationEventListener, n *broker.NatsMessageBroker, cfg *config.WarehouseConfig) *ReservationEventRouter {
	return &ReservationEventRouter{reservationListener, n, cfg, restoreFactory.Build()}
}

func (r *ReservationEventRouter) Setup(ctx context.Context) error {
	// register stream message handlers
	err := r.broker.RegisterJsHandler(
		ctx, r.restore, stream.ReservationEventStreamConfig, r.reservationListener.ListenReservationEvent,
		broker.WithSubjectFilter(fmt.Sprintf("reservation.%s", r.cfg.ID)),
	)
	if err != nil {
		return err
	}

	// wait restoring of the state before starting the server
	r.restore.Wait()

	return nil
}
