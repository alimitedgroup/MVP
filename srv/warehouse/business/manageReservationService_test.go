package business

import (
	"context"
	"testing"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
)

func TestManageReservationServiceApplyReservationEvent(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)

	cfg := config.WarehouseConfig{ID: "1"}

	createStockUpdatePortMock := NewMockICreateStockUpdatePort(ctrl)
	getStockPortMock := NewMockIGetStockPort(ctrl)
	storeReservationEventPort := NewMockIStoreReservationEventPort(ctrl)
	getReservationPort := NewMockIGetReservationPort(ctrl)

	applyReservationEventPort := NewMockIApplyReservationEventPort(ctrl)
	applyReservationEventPort.EXPECT().ApplyReservationEvent(gomock.Any()).Return(nil)

	idempotentPortMock := NewMockIIdempotentPort(ctrl)
	idempotentPortMock.EXPECT().IsAlreadyProcessed(gomock.Any()).Return(false)

	app := fx.New(
		fx.Supply(&cfg),
		fx.Supply(fx.Annotate(createStockUpdatePortMock, fx.As(new(port.ICreateStockUpdatePort)))),
		fx.Supply(fx.Annotate(getStockPortMock, fx.As(new(port.IGetStockPort)))),
		fx.Supply(fx.Annotate(storeReservationEventPort, fx.As(new(port.IStoreReservationEventPort)))),
		fx.Supply(fx.Annotate(idempotentPortMock, fx.As(new(port.IIdempotentPort)))),
		fx.Supply(fx.Annotate(applyReservationEventPort, fx.As(new(port.IApplyReservationEventPort)))),
		fx.Supply(fx.Annotate(getReservationPort, fx.As(new(port.IGetReservationPort)))),
		fx.Provide(NewManageReservationService),
		fx.Invoke(func(lc fx.Lifecycle, service *ManageReservationService) {
			lc.Append(fx.Hook{OnStart: func(ctx context.Context) error {
				applyCmd := port.ApplyReservationEventCmd{}
				err := service.ApplyReservationEvent(applyCmd)
				require.NoError(t, err)

				return nil
			}})
		}),
	)

	err := app.Start(ctx)
	require.NoError(t, err)

	defer func() {
		err := app.Stop(ctx)
		require.NoError(t, err)
	}()

}
