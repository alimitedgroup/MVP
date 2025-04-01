package adapterin

import (
	"context"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/notification/portin"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
)

type startResult struct {
	queryResult  *MockQueryRules
	stockUpdates *MockStockUpdates
	nc           *nats.Conn
	js           jetstream.JetStream
}

//go:generate go run go.uber.org/mock/mockgen@latest -destination business_query_rules_mock.go -package adapterin github.com/alimitedgroup/MVP/srv/notification/portin QueryRules
//go:generate go run go.uber.org/mock/mockgen@latest -destination business_stock_updates_mock.go -package adapterin github.com/alimitedgroup/MVP/srv/notification/portin StockUpdates

func start(t *testing.T) startResult {
	ctrl := gomock.NewController(t)
	queryResult := NewMockQueryRules(ctrl)
	stockUpdatesMock := NewMockStockUpdates(ctrl)

	nc, _ := broker.NewInProcessNATSServer(t)

	js, err := jetstream.New(nc)
	require.NoError(t, err)

	app := fx.New(
		Module,
		lib.ModuleTest,
		fx.Supply(
			fx.Annotate(queryResult, fx.As(new(portin.QueryRules))),
			fx.Annotate(stockUpdatesMock, fx.As(new(portin.StockUpdates))),
			nc, t,
		),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = app.Start(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := app.Stop(ctx)
		require.NoError(t, err)
	})

	return startResult{
		queryResult:  queryResult,
		stockUpdates: stockUpdatesMock,
		nc:           nc,
		js:           js,
	}
}
