package controller

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestHealthCheckController(t *testing.T) {
	ctx := t.Context()

	ns, _ := broker.NewInProcessNATSServer(t)

	app := fx.New(
		fx.Supply(ns),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(NewHealthCheckController),
		fx.Provide(NewHealthCheckRouter),
		fx.Invoke(func(lc fx.Lifecycle, r *HealthCheckRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					require.NoError(t, err)

					resp, err := ns.Request("order.ping", []byte{}, 1*time.Second)
					require.NoError(t, err)

					var respDto response.HealthCheckResponseDTO
					err = json.Unmarshal(resp.Data, &respDto)
					require.NoError(t, err)

					require.Empty(t, respDto.Error)
					require.Equal(t, respDto.Message, "pong")

					return nil
				},
			})
		}),
	)

	err := app.Start(ctx)
	require.NoError(t, err)

	defer func() {
		err = app.Stop(ctx)
		require.NoError(t, err)
	}()
}
