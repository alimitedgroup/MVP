package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	brokerRouter "github.com/alimitedgroup/MVP/srv/warehouse/adapter/controller"
	"github.com/magiconair/properties/assert"
	"go.uber.org/fx"
)

func TestRunWithBadConfigParams(t *testing.T) {
	ctx := context.Background()
	cfg := broker.BrokerConfig{
		Url: "nats://localhost:-100",
	}

	app := fx.New(
		fx.Supply(&cfg),
		fx.Supply(brokerRouter.BrokerRoutes{}),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Invoke(func(config *broker.BrokerConfig, broker *broker.NatsMessageBroker) {
			assert.Equal(t, cfg.Url, config.Url)
			assert.Equal(t, broker, nil)
		}),
		fx.Invoke(RunLifeCycle),
	)

	err := app.Start(ctx)
	assert.Equal(t, err != nil, true, fmt.Sprintf("expected error on connecting to NATS with port address %v", cfg.Url))
}
