package listener

import (
	"context"
	"fmt"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/nats-io/nats.go/jetstream"
)

type OrderUpdateRouter struct {
	reservationListener *OrderUpdateListener
	broker              *broker.NatsMessageBroker
	restore             broker.IRestoreStreamControl
	cfg                 *config.WarehouseConfig
}

func NewOrderUpdateRouter(restoreFactory broker.IRestoreStreamControlFactory, reservationListener *OrderUpdateListener, broker *broker.NatsMessageBroker, cfg *config.WarehouseConfig) *OrderUpdateRouter {
	return &OrderUpdateRouter{reservationListener, broker, restoreFactory.Build(), cfg}
}

func (r *OrderUpdateRouter) Setup(ctx context.Context) error {
	// register stream message handlers
	orderUpdateConsumerConfig := jetstream.ConsumerConfig{
		Durable: fmt.Sprintf("order-update-warehouse-%s", r.cfg.ID),
	}

	err := r.broker.RegisterJsWithConsumerGroup(ctx, stream.OrderUpdateStreamConfig, orderUpdateConsumerConfig, r.reservationListener.ListenOrderUpdate)
	if err != nil {
		return err
	}

	transferUpdateConsumerConfig := jetstream.ConsumerConfig{
		Durable: fmt.Sprintf("transfer-update-warehouse-%s", r.cfg.ID),
	}
	err = r.broker.RegisterJsWithConsumerGroup(ctx, stream.TransferUpdateStreamConfig, transferUpdateConsumerConfig, r.reservationListener.ListenTransferUpdate)
	if err != nil {
		return err
	}

	// wait restoring of the state before starting the server
	r.restore.Wait()

	return nil
}
