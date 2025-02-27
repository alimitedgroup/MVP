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
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type RunParams struct {
	fx.In

	ServerConfig *APIConfig
	HttpHandler  *lib.HTTPHandler
	ApiRoutes    apiRouter.APIRoutes
	BrokerRoutes brokerRouter.BrokerRoutes
}

func Run(p RunParams) {
	p.ApiRoutes.Setup()
	p.BrokerRoutes.Setup()

	err := p.HttpHandler.Engine.Run(fmt.Sprintf(":%d", p.ServerConfig.Port))
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
