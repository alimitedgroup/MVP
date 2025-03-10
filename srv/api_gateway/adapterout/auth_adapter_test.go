package adapterout

import (
	"context"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portout"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"testing"
	"time"
)

func createAuthAdapter(t *testing.T, nc *nats.Conn) portout.AuthenticationPortOut {
	brk, err := broker.NewNatsMessageBroker(nc)
	require.NoError(t, err)
	return NewAuthenticationAdapter(brk)
}

func startAuthMock(t *testing.T, nc *nats.Conn, issuer string) func() {
	app := fx.New(
		fx.Supply(issuer),
		fx.Supply(nc),
		fx.Supply(DefaultUsers),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(NewAuthMock),
		fx.Invoke(StartAuthMock),
	)

	// Start microservice, with timeout of 1 second
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	t.Cleanup(cancel)
	require.NoError(t, app.Start(ctx))

	// Stop microservice, with timeout of 1 second
	stopFunc := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		t.Cleanup(cancel)
		require.NoError(t, app.Stop(ctx))
	}

	t.Cleanup(stopFunc)
	return stopFunc
}

func TestLogin(t *testing.T) {
	nc := broker.NewInProcessNATSServer(t)

	aa := createAuthAdapter(t, nc)
	startAuthMock(t, nc, "test")

	token, err := aa.GetToken("admin")
	require.NoError(t, err)
	res, err := aa.VerifyToken(token)
	require.NoError(t, err)

	username, err := aa.GetUsername(res)
	require.NoError(t, err)
	require.Equal(t, "admin", username)

	role, err := aa.GetRole(res)
	require.NoError(t, err)
	require.Equal(t, types.RoleGlobalAdmin, role)
}

func TestRevocation(t *testing.T) {
	nc := broker.NewInProcessNATSServer(t)

	aa := createAuthAdapter(t, nc)
	stopEarly := startAuthMock(t, nc, "test")

	token, err := aa.GetToken("admin")
	require.NoError(t, err)
	_, err = aa.VerifyToken(token)
	require.NoError(t, err)

	stopEarly()
	startAuthMock(t, nc, "test")

	res, err := aa.VerifyToken(token)
	require.ErrorIs(t, err, portout.ErrTokenInvalid)
	require.Nil(t, res)
}
