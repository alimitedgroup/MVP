package controller

import (
	"context"
	"testing"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	serviceportin "github.com/alimitedgroup/MVP/srv/authenticator/service/portIn"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	gomock "go.uber.org/mock/gomock"
)

func Test_Router(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)

	ns, _ := broker.NewInProcessNATSServer(t)

	app := fx.New(
		Module,
		fx.Supply(ns),
		fx.Supply(ctrl),
		fx.Provide(
			fx.Annotate(NewFakeService,
				fx.As(new(serviceportin.IGetTokenUseCase)),
			)),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(broker.NewRestoreStreamControl),
		fx.Invoke(func(lc fx.Lifecycle, r *ControllerRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					require.NoError(t, err)
					return nil
				},
			})
		}),
	)

	err := app.Start(ctx)
	require.NoError(t, err)
}
