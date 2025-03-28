package controller

import (
	"context"
	"github.com/alimitedgroup/MVP/common/lib"
	"testing"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination mock_controller.go -package controller github.com/alimitedgroup/MVP/srv/order/business/port ICreateTransferUseCase,IGetTransferUseCase,ICreateOrderUseCase,IGetOrderUseCase

func TestRouter(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)

	ns, _ := broker.NewInProcessNATSServer(t)

	app := fx.New(
		Module,
		lib.ModuleTest,
		fx.Supply(ns, t, ctrl),
		fx.Provide(fx.Annotate(NewMockICreateTransferUseCase, fx.As(new(port.ICreateTransferUseCase)))),
		fx.Provide(fx.Annotate(NewMockIGetTransferUseCase, fx.As(new(port.IGetTransferUseCase)))),
		fx.Provide(fx.Annotate(NewMockICreateOrderUseCase, fx.As(new(port.ICreateOrderUseCase)))),
		fx.Provide(fx.Annotate(NewMockIGetOrderUseCase, fx.As(new(port.IGetOrderUseCase)))),
		fx.Invoke(func(lc fx.Lifecycle, r *ControllerRoutes) {
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
