package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/alimitedgroup/MVP/common/lib"
	apiRouter "github.com/alimitedgroup/MVP/srv/api_gateway/api/router"
	brokerRouter "github.com/alimitedgroup/MVP/srv/api_gateway/channel/router"
	"github.com/magiconair/properties/assert"
	"go.uber.org/fx"
)

func TestRunWithBadConfigParams(t *testing.T) {
	ctx := context.Background()

	config := APIConfig{
		Host: "localhost",
		Port: -100,
	}

	app := fx.New(
		fx.Provide(lib.NewHTTPHandler),
		fx.Supply(apiRouter.APIRoutes{}),
		fx.Supply(brokerRouter.BrokerRoutes{}),
		fx.Supply(&config),
		fx.Invoke(func(p RunParams) {
			assert.Equal(t, p.ServerConfig.Host, config.Host)
			assert.Equal(t, p.ServerConfig.Port, config.Port)
		}),
		fx.Invoke(RunLifeCycle),
	)

	err := app.Start(ctx)
	assert.Equal(t, err != nil, true, fmt.Sprintf("expected error on listening to port %d", config.Port))
}
