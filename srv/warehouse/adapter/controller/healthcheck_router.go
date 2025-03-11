package controller

import (
	"context"
	"fmt"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
)

type HealthCheckRouter struct {
	config                *config.WarehouseConfig
	healthCheckController *HealthCheckController
	broker                *broker.NatsMessageBroker
}

func NewHealthCheckRouter(config *config.WarehouseConfig, healthCheckController *HealthCheckController, broker *broker.NatsMessageBroker) *HealthCheckRouter {
	return &HealthCheckRouter{config, healthCheckController, broker}
}

func (r *HealthCheckRouter) Setup(ctx context.Context) error {
	// register request/reply handlers
	err := r.broker.RegisterRequest(ctx, broker.Subject(fmt.Sprintf("warehouse.%s.ping", r.config.ID)), broker.NoQueue, r.healthCheckController.PingHandler)
	if err != nil {
		return err
	}

	return nil
}
