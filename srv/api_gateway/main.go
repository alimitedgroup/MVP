package main

import (
	"fmt"
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/common/lib/broker"
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

func main() {
	config := loadConfig()

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", 8080))
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
		fx.Provide(broker.NewNatsConn),
		fx.Provide(adapterin.NewListener),
	)

	app.Run()
}
