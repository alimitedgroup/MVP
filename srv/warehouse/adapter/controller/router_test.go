package controller

import (
	"context"
	"testing"

	"github.com/alimitedgroup/MVP/common/lib"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination mock_controller.go -package controller github.com/alimitedgroup/MVP/srv/warehouse/business/port IAddStockUseCase,IRemoveStockUseCase,ICreateReservationUseCase

func TestRouter(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)

	ns, _ := broker.NewInProcessNATSServer(t)
	cfg := config.WarehouseConfig{ID: "1"}

	app := fx.New(
		Module,
		lib.ModuleTest,
		fx.Supply(ns, t, &cfg, ctrl),
		fx.Provide(fx.Annotate(NewMockIAddStockUseCase, fx.As(new(port.IAddStockUseCase)))),
		fx.Provide(fx.Annotate(NewMockIRemoveStockUseCase, fx.As(new(port.IRemoveStockUseCase)))),
		fx.Provide(observability.TestLogger),
		fx.Provide(observability.TestMeter),
		fx.Provide(fx.Annotate(NewMockICreateReservationUseCase, fx.As(new(port.ICreateReservationUseCase)))),
		fx.Invoke(func(lc fx.Lifecycle, r *BrokerRoutes) {
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
