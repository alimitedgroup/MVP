package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/srv/api_gateway/api"
	apiRouter "github.com/alimitedgroup/MVP/srv/api_gateway/api/router"
	"github.com/alimitedgroup/MVP/srv/api_gateway/channel"
	brokerRouter "github.com/alimitedgroup/MVP/srv/api_gateway/channel/router"
	"go.uber.org/fx"
)

type APIConfig struct {
	Host string
	Port int
}

func Run(serverConfig *APIConfig, h *lib.HTTPHandler, apiRoutes apiRouter.APIRoutes, brokerRoutes brokerRouter.BrokerRoutes) {
	apiRoutes.Setup()
	brokerRoutes.Setup()

	err := h.Engine.Run(fmt.Sprintf(":%d", serverConfig.Port))
	if err != nil {
		log.Fatal("error running the Gin HTTP engine\n", err)
	}
}

var Modules = fx.Options(
	lib.Module,
	api.Module,
	channel.Module,
)

func main() {
	_ = context.Background()

	config := loadConfig()

	opts := fx.Options(
		config,
		Modules,
	)

	app := fx.New(
		opts,
		fx.Invoke(Run),
	)

	app.Run()
}
