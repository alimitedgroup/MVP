package main

import (
	"context"
	"fmt"
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/srv/api_gateway/adapterin"
	"github.com/alimitedgroup/MVP/srv/api_gateway/adapterout"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business"
	"go.uber.org/fx"
	"log"
	"net"
)

type APIConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type RunParams struct {
	fx.In

	ServerConfig *APIConfig
	HttpHandler  *adapterin.HTTPHandler
}

func Run(p RunParams) error {
	err := p.HttpHandler.Engine.Run(fmt.Sprintf(":%d", p.ServerConfig.Port))
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

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", "localhost", 8080))
	if err != nil {
		log.Fatal("Invalid TCP address: ", err)
	}

	app := fx.New(
		config,
		lib.Module,
		business.Module,
		adapterout.Module,
		adapterin.Module,
		fx.Supply(addr),
		fx.Provide(adapterin.NewListener),
		fx.Invoke(RunLifeCycle),
	)

	err = app.Start(ctx)
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
