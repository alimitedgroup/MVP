package controller

import (
	"context"
	"github.com/alimitedgroup/MVP/common/lib"
	"testing"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	serviceportin "github.com/alimitedgroup/MVP/srv/authenticator/service/portIn"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
)

func Test_Router(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)

	ns, _ := broker.NewInProcessNATSServer(t)

	app := fx.New(
		Module,
		lib.ModuleTest,
		fx.Supply(ns, t, ctrl),
		fx.Provide(fx.Annotate(NewFakeService, fx.As(new(serviceportin.IGetTokenUseCase)))),
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
