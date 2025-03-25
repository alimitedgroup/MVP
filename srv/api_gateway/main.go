package main

import (
	"fmt"
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/srv/api_gateway/adapterin"
	"github.com/alimitedgroup/MVP/srv/api_gateway/adapterout"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business"
	prettyconsole "github.com/thessem/zap-prettyconsole"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
)

type APIConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func main() {
	config := loadConfig()

	addrStr := fmt.Sprintf("%s:%d", "0.0.0.0", 8080)
	addr, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		prettyconsole.NewLogger(zap.DebugLevel).Fatal(
			"Failed to bind to TCP address",
			zap.Error(err),
			zap.String("addr", addrStr),
		)
	}

	app := fx.New(
		config,
		lib.Module,
		business.Module,
		adapterout.Module,
		adapterin.Module,
		observability.Module,
		fx.Supply(addr),
		fx.Provide(broker.NewNatsConn),
		fx.Provide(adapterin.NewListener),
	)

	app.Run()
}
