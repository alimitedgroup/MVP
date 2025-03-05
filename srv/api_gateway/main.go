package main

import (
	"context"
	"fmt"
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/srv/api_gateway/adapterout"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business"
	"github.com/alimitedgroup/MVP/srv/api_gateway/controller"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portout"
	"go.uber.org/fx"
	"log"
)

type APIConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type RunParams struct {
	fx.In

	ServerConfig *APIConfig
	HttpHandler  *lib.HTTPHandler
}

func Run(p RunParams) error {
	var err error

	err = p.HttpHandler.Engine.Run(fmt.Sprintf(":%d", p.ServerConfig.Port))
	if err != nil {
		return err
	}

	return nil
}

func RunLifeCycle(lc fx.Lifecycle, p RunParams) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := Run(p)
			return err
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

func main() {
	ctx := context.Background()

	config := loadConfig()

	app := fx.New(
		lib.Module,
		business.Module,
		fx.Provide(
			fx.Annotate(adapterout.NewAuthenticationAdapter, fx.As(new(portout.AuthenticationPortOut))),
		),
		controller.Module,
		config,
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
