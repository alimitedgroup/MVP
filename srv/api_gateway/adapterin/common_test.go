package adapterin

import (
	"context"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
	"net"
	"testing"
	"time"
)

type startResult struct {
	base       string
	auth       *MockAuth
	warehouses *MockWarehouses
}

//go:generate go run go.uber.org/mock/mockgen@latest -destination business_auth_mock.go -package adapterin github.com/alimitedgroup/MVP/srv/api_gateway/portin Auth
//go:generate go run go.uber.org/mock/mockgen@latest -destination business_warehouses_mock.go -package adapterin github.com/alimitedgroup/MVP/srv/api_gateway/portin Warehouses

// start starts the application with a mock business login,
// and returns it, along with the base url that can be used to send requests
func start(t *testing.T) startResult {
	ctrl := gomock.NewController(t)
	mock := NewMockAuth(ctrl)
	wMock := NewMockWarehouses(ctrl)

	nc, _ := broker.NewInProcessNATSServer(t)

	addr, err := net.ResolveTCPAddr("tcp", ":0")
	require.NoError(t, err)
	ln, err := net.ListenTCP("tcp", addr)
	require.NoError(t, err)

	app := fx.New(
		Module,
		fx.Supply(zaptest.NewLogger(t)),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Supply(
			fx.Annotate(mock, fx.As(new(portin.Auth))),
			fx.Annotate(wMock, fx.As(new(portin.Warehouses))),
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
