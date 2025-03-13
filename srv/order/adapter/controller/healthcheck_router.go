package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib/broker"
)

type HealthCheckRouter struct {
	healthCheckController *HealthCheckController
	broker                *broker.NatsMessageBroker
}

func NewHealthCheckRouter(healthCheckController *HealthCheckController, broker *broker.NatsMessageBroker) *HealthCheckRouter {
	return &HealthCheckRouter{healthCheckController, broker}
}

func (r *HealthCheckRouter) Setup(ctx context.Context) error {
	// register request/reply handlers
	err := r.broker.RegisterRequest(ctx, broker.Subject("order.ping"), broker.NoQueue, r.healthCheckController.PingHandler)
	if err != nil {
		return err
	}

	return nil
}
