package adapterin

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
)

type startResult struct {
	base          string
	auth          *MockAuth
	warehouses    *MockWarehouses
	order         *MockOrder
	notifications *MockNotifications
}

//go:generate go run go.uber.org/mock/mockgen@latest -destination business_auth_mock.go -package adapterin github.com/alimitedgroup/MVP/srv/api_gateway/portin Auth
//go:generate go run go.uber.org/mock/mockgen@latest -destination business_warehouses_mock.go -package adapterin github.com/alimitedgroup/MVP/srv/api_gateway/portin Warehouses
//go:generate go run go.uber.org/mock/mockgen@latest -destination business_order_mock.go -package adapterin github.com/alimitedgroup/MVP/srv/api_gateway/portin Order
//go:generate go run go.uber.org/mock/mockgen@latest -destination business_notifications_mock.go -package adapterin github.com/alimitedgroup/MVP/srv/api_gateway/portin Notifications

// start starts the application with a mock business login,
// and returns it, along with the base url that can be used to send requests
func start(t *testing.T) startResult {
	ctrl := gomock.NewController(t)
	mock := NewMockAuth(ctrl)
	wMock := NewMockWarehouses(ctrl)
	orderMock := NewMockOrder(ctrl)
	notificationsMock := NewMockNotifications(ctrl)

	nc, _ := broker.NewInProcessNATSServer(t)

	addr, err := net.ResolveTCPAddr("tcp", ":0")
	require.NoError(t, err)
	ln, err := net.ListenTCP("tcp", addr)
	require.NoError(t, err)

	app := fx.New(
		ModuleTest,
		lib.ModuleTest,
		fx.Supply(
			fx.Annotate(mock, fx.As(new(portin.Auth))),
			fx.Annotate(wMock, fx.As(new(portin.Warehouses))),
			fx.Annotate(orderMock, fx.As(new(portin.Order))),
			fx.Annotate(notificationsMock, fx.As(new(portin.Notifications))),
			ln, nc, t,
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
		auth:          mock,
		warehouses:    wMock,
		order:         orderMock,
		notifications: notificationsMock,
		base:          "http://" + ln.Addr().String(),
	}
}
