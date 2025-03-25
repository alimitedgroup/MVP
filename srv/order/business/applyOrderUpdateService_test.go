package business

import (
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	gomock "go.uber.org/mock/gomock"
)

type applyOrderUpdateServiceMockSuite struct {
	applyOrderUpdatePort    *MockIApplyOrderUpdatePort
	applyTransferUpdatePort *MockIApplyTransferUpdatePort
}

func newApplyOrderUpdateServiceMockSuite(t *testing.T) *applyOrderUpdateServiceMockSuite {
	ctrl := gomock.NewController(t)

	return &applyOrderUpdateServiceMockSuite{
		applyOrderUpdatePort:    NewMockIApplyOrderUpdatePort(ctrl),
		applyTransferUpdatePort: NewMockIApplyTransferUpdatePort(ctrl),
	}
}

func runTestapplyOrderUpdateService(t *testing.T, build func(*applyOrderUpdateServiceMockSuite), buildOptions func() fx.Option, runLifeCycle func() interface{}) {
	ctx := t.Context()
	suite := newApplyOrderUpdateServiceMockSuite(t)

	build(suite)

	app := fx.New(
		fx.Supply(fx.Annotate(suite.applyOrderUpdatePort, fx.As(new(port.IApplyOrderUpdatePort)))),
		fx.Supply(fx.Annotate(suite.applyTransferUpdatePort, fx.As(new(port.IApplyTransferUpdatePort)))),
		fx.Provide(NewApplyOrderUpdateService),
		fx.Invoke(runLifeCycle()),
		buildOptions(),
	)

	err := app.Start(ctx)
	require.NoError(t, err)

	defer func() {
		err := app.Stop(ctx)
		require.NoError(t, err)
	}()
}

func TestApplyOrderUpdateServiceOrder(t *testing.T) {
	ctx := t.Context()
	runTestapplyOrderUpdateService(t,
		func(suite *applyOrderUpdateServiceMockSuite) {
			suite.applyOrderUpdatePort.EXPECT().ApplyOrderUpdate(gomock.Any())
		},
		func() fx.Option { return fx.Options() },
		func() interface{} {
			return func(service *ApplyOrderUpdateService) {
				cmd := port.OrderUpdateCmd{
					ID:           "1",
					Status:       "Created",
					Name:         "test",
					FullName:     "test test",
					Address:      "via roma 1",
					CreationTime: time.Now().UnixMilli(),
					UpdateTime:   time.Now().UnixMilli(),
					Reservations: []string{},
					Goods: []port.OrderUpdateGood{
						{
							GoodID:   "1",
							Quantity: 1,
						},
					},
				}
				service.ApplyOrderUpdate(ctx, cmd)
			}
		},
	)
}

func TestApplyOrderUpdateServiceTransfer(t *testing.T) {
	ctx := t.Context()
	runTestapplyOrderUpdateService(t,
		func(suite *applyOrderUpdateServiceMockSuite) {
			suite.applyTransferUpdatePort.EXPECT().ApplyTransferUpdate(gomock.Any())
		},
		func() fx.Option { return fx.Options() },
		func() interface{} {
			return func(service *ApplyOrderUpdateService) {
				cmd := port.TransferUpdateCmd{
					ID:            "1",
					Status:        "Created",
					ReceiverID:    "1",
					SenderID:      "2",
					CreationTime:  time.Now().UnixMilli(),
					UpdateTime:    time.Now().UnixMilli(),
					ReservationID: "",
					Goods: []port.TransferUpdateGood{
						{
							GoodID:   "1",
							Quantity: 1,
						},
					},
				}
				service.ApplyTransferUpdate(ctx, cmd)
			}
		},
	)
}
