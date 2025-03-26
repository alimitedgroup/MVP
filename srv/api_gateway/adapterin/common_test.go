package adapterin

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

type startResult struct {
	base       string
	auth       *MockAuth
	warehouses *MockWarehouses
}

//go:generate go run go.uber.org/mock/mockgen@latest -destination business_auth_mock.go -package adapterin github.com/alimitedgroup/MVP/srv/api_gateway/portin Auth
//go:generate go run go.uber.org/mock/mockgen@latest -destination business_warehouses_mock.go -package adapterin github.com/alimitedgroup/MVP/srv/api_gateway/portin Warehouses
//go:generate go run go.uber.org/mock/mockgen@latest -destination business_order_mock.go -package adapterin github.com/alimitedgroup/MVP/srv/api_gateway/portin Order

// start starts the application with a mock business login,
// and returns it, along with the base url that can be used to send requests
func start(t *testing.T) startResult {
	ctrl := gomock.NewController(t)
	mock := NewMockAuth(ctrl)
	wMock := NewMockWarehouses(ctrl)
	orderMock := NewMockOrder(ctrl)

	nc, _ := broker.NewInProcessNATSServer(t)

	addr, err := net.ResolveTCPAddr("tcp", ":0")
	require.NoError(t, err)
	ln, err := net.ListenTCP("tcp", addr)
	require.NoError(t, err)

	app := fx.New(
		Module,
		fx.Supply(zaptest.NewLogger(t)),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(func() metric.Meter { return otel.Meter("test") }),
		fx.Supply(
			fx.Annotate(mock, fx.As(new(portin.Auth))),
			fx.Annotate(wMock, fx.As(new(portin.Warehouses))),
			fx.Annotate(orderMock, fx.As(new(portin.Order))),
			ln, nc,
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
		auth:       mock,
		warehouses: wMock,
		base:       "http://" + ln.Addr().String(),
	}
}
