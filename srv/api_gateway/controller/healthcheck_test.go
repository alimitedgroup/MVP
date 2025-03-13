package controller

import (
	"context"
	"encoding/json"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
	"net"
	"net/http"
	"testing"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination business_auth_mock.go -package controller github.com/alimitedgroup/MVP/srv/api_gateway/portin Auth

func start(t *testing.T) string {
	ctrl := gomock.NewController(t)
	mock := NewMockAuth(ctrl)

	nc, _ := broker.NewInProcessNATSServer(t)

	addr, err := net.ResolveTCPAddr("tcp", ":0")
	require.NoError(t, err)
	ln, err := net.ListenTCP("tcp", addr)
	require.NoError(t, err)

	app := fx.New(
		Module,
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Supply(
			fx.Annotate(mock, fx.As(new(portin.Auth))),
			ln, nc,
		),
	)

	err = app.Start(context.Background())
	require.NoError(t, err)

	t.Cleanup(func() {
		err := app.Stop(context.Background())
		require.NoError(t, err)
	})

	return "http://" + ln.Addr().String()
}

func TestHealthCheck(t *testing.T) {
	base := start(t)
	client := &http.Client{}

	resp, err := client.Get(base + "/api/v1/ping")
	require.NoError(t, err)

	var respbody response.ResponseDTO[string]
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)

	require.Equal(t, "pong", respbody.Message)
	require.Zero(t, respbody.Error)
}
