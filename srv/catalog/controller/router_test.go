package controller

import (
	"context"
	"testing"

	"go.uber.org/zap/zaptest"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	serviceportin "github.com/alimitedgroup/MVP/srv/catalog/service/portin"
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
		fx.Supply(ns),
		fx.Supply(ctrl),
		fx.Supply(zaptest.NewLogger(t)),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(broker.NewRestoreStreamControl),
		fx.Provide(observability.TestMeter),
		fx.Provide(
			fx.Annotate(NewFakeControllerUC,
				fx.As(new(serviceportin.IGetGoodsInfoUseCase)),
				fx.As(new(serviceportin.IGetGoodsQuantityUseCase)),
				fx.As(new(serviceportin.IGetWarehousesUseCase)),
				fx.As(new(serviceportin.ISetMultipleGoodsQuantityUseCase)),
				fx.As(new(serviceportin.IUpdateGoodDataUseCase)),
			),
		),
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
