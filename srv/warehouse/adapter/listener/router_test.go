package listener

import (
	"context"
	"testing"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination=mock_listener.go -package listener github.com/alimitedgroup/MVP/srv/warehouse/business/port IApplyReservationUseCase,IApplyCatalogUpdateUseCase,IConfirmOrderUseCase,IConfirmTransferUseCase,IApplyStockUpdateUseCase

func TestRouter(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)

	ns, _ := broker.NewInProcessNATSServer(t)
	cfg := config.WarehouseConfig{ID: "1"}

	app := fx.New(
		Module,
		fx.Supply(&cfg),
		fx.Supply(ns),
		fx.Supply(ctrl),
		fx.Provide(fx.Annotate(NewMockIApplyReservationUseCase, fx.As(new(port.IApplyReservationUseCase)))),
		fx.Provide(fx.Annotate(NewMockIApplyCatalogUpdateUseCase, fx.As(new(port.IApplyCatalogUpdateUseCase)))),
		fx.Provide(fx.Annotate(NewMockIConfirmOrderUseCase, fx.As(new(port.IConfirmOrderUseCase)))),
		fx.Provide(fx.Annotate(NewMockIConfirmTransferUseCase, fx.As(new(port.IConfirmTransferUseCase)))),
		fx.Provide(fx.Annotate(NewMockIApplyStockUpdateUseCase, fx.As(new(port.IApplyStockUpdateUseCase)))),
		fx.Supply(zaptest.NewLogger(t)),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(fx.Annotate(broker.NewRestoreStreamControlFactory, fx.As(new(broker.IRestoreStreamControlFactory)))),
		fx.Invoke(func(lc fx.Lifecycle, r *ListenerRoutes) {
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
