package listener

import (
	"context"
	"testing"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination mock_listener.go -package listener github.com/alimitedgroup/MVP/srv/order/business/port IApplyOrderUpdateUseCase,IApplyTransferUpdateUseCase,IContactWarehousesUseCase,IApplyStockUpdateUseCase

func TestRouter(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)

	ns, _ := broker.NewInProcessNATSServer(t)

	app := fx.New(
		Module,
		fx.Supply(ns),
		fx.Supply(ctrl),
		fx.Provide(fx.Annotate(NewMockIApplyOrderUpdateUseCase, fx.As(new(port.IApplyOrderUpdateUseCase)))),
		fx.Provide(fx.Annotate(NewMockIApplyTransferUpdateUseCase, fx.As(new(port.IApplyTransferUpdateUseCase)))),
		fx.Provide(fx.Annotate(NewMockIApplyStockUpdateUseCase, fx.As(new(port.IApplyStockUpdateUseCase)))),
		fx.Provide(fx.Annotate(NewMockIContactWarehousesUseCase, fx.As(new(port.IContactWarehousesUseCase)))),
		fx.Provide(fx.Annotate(broker.NewRestoreStreamControlFactory, fx.As(new(broker.IRestoreStreamControlFactory)))),
		fx.Provide(broker.NewNatsMessageBroker),
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
