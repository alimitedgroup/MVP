package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/srv/api_gateway/api"
	apiController "github.com/alimitedgroup/MVP/srv/api_gateway/api/controller"
	"github.com/alimitedgroup/MVP/srv/api_gateway/channel"
	brokerController "github.com/alimitedgroup/MVP/srv/api_gateway/channel/controller"
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
	ApiRoutes    apiController.APIRoutes
	BrokerRoutes brokerController.BrokerRoutes
}

func Run(ctx context.Context, p RunParams) error {
	var err error

	err = p.BrokerRoutes.Setup(ctx)
	if err != nil {
		return err
	}

	p.ApiRoutes.Setup(ctx)

	err = p.HttpHandler.Engine.Run(fmt.Sprintf(":%d", p.ServerConfig.Port))
	if err != nil {
		return err
	}

	return nil
}

func RunLifeCycle(lc fx.Lifecycle, p RunParams) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := Run(ctx, p)
			return err
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

var Modules = fx.Options(
	lib.Module,
	api.Module,
	channel.Module,
)

func main() {
	ctx := context.Background()

	config := loadConfig()

	opts := fx.Options(
		config,
		Modules,
	)

	app := fx.New(
		opts,
		fx.Invoke(RunLifeCycle),
	)

	err := app.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = app.Stop(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

}
