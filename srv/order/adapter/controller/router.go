package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
)

type ControllerRoutes struct {
	broker *broker.NatsMessageBroker
	routes []lib.BrokerRoute
}

func NewBrokerRoutes(stockRouter *OrderRouter, healthCheckRouter *HealthCheckRouter, broker *broker.NatsMessageBroker) ControllerRoutes {
	return ControllerRoutes{
		routes: []lib.BrokerRoute{stockRouter, healthCheckRouter},
		broker: broker,
	}
}

func (r ControllerRoutes) Setup(ctx context.Context) error {
	_, err := r.broker.Js.CreateOrUpdateStream(ctx, stream.OrderUpdateStreamConfig)
	if err != nil {
		return err
	}

	for _, v := range r.routes {
		err := v.Setup(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
