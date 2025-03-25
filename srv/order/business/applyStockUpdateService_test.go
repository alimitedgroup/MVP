package business

import (
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	gomock "go.uber.org/mock/gomock"
)

type applyStockUpdateServiceMockSuite struct {
	applyStockUpdatePort      *MockIApplyStockUpdatePort
	applyOrderUpdatePort      *MockIApplyOrderUpdatePort
	getOrderPort              *MockIGetOrderPort
	getTransferPort           *MockIGetTransferPort
	applyTransferUpdatePort   *MockIApplyTransferUpdatePort
	setCompleteTransferPort   *MockISetCompleteTransferPort
	setCompletedWarehousePort *MockISetCompletedWarehouseOrderPort
}

func newApplyStockUpdateServiceMockSuite(t *testing.T) *applyStockUpdateServiceMockSuite {
	ctrl := gomock.NewController(t)

	return &applyStockUpdateServiceMockSuite{
		applyStockUpdatePort:      NewMockIApplyStockUpdatePort(ctrl),
		applyOrderUpdatePort:      NewMockIApplyOrderUpdatePort(ctrl),
		getOrderPort:              NewMockIGetOrderPort(ctrl),
		getTransferPort:           NewMockIGetTransferPort(ctrl),
		applyTransferUpdatePort:   NewMockIApplyTransferUpdatePort(ctrl),
		setCompleteTransferPort:   NewMockISetCompleteTransferPort(ctrl),
		setCompletedWarehousePort: NewMockISetCompletedWarehouseOrderPort(ctrl),
	}
}

func runTestApplyStockUpdateService(t *testing.T, build func(*applyStockUpdateServiceMockSuite), buildOptions func() fx.Option, runLifeCycle func() interface{}) {
	ctx := t.Context()
	suite := newApplyStockUpdateServiceMockSuite(t)

	build(suite)

	app := fx.New(
		fx.Supply(fx.Annotate(suite.applyStockUpdatePort, fx.As(new(port.IApplyStockUpdatePort)))),
		fx.Supply(fx.Annotate(suite.applyOrderUpdatePort, fx.As(new(port.IApplyOrderUpdatePort)))),
		fx.Supply(fx.Annotate(suite.getOrderPort, fx.As(new(port.IGetOrderPort)))),
		fx.Supply(fx.Annotate(suite.getTransferPort, fx.As(new(port.IGetTransferPort)))),
		fx.Supply(fx.Annotate(suite.applyTransferUpdatePort, fx.As(new(port.IApplyTransferUpdatePort)))),
		fx.Supply(fx.Annotate(suite.setCompleteTransferPort, fx.As(new(port.ISetCompleteTransferPort)))),
		fx.Supply(fx.Annotate(suite.setCompletedWarehousePort, fx.As(new(port.ISetCompletedWarehouseOrderPort)))),
		fx.Provide(NewApplyStockUpdateService),
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

func TestApplyStockUpdateServiceStock(t *testing.T) {
	ctx := t.Context()
	runTestApplyStockUpdateService(t,
		func(suite *applyStockUpdateServiceMockSuite) {
			now := time.Now().UnixMilli()
			suite.applyStockUpdatePort.EXPECT().ApplyStockUpdate(gomock.Any())
			suite.getOrderPort.EXPECT().GetOrder(gomock.Any()).Return(model.Order{
				ID:           "1",
				Status:       "Filled",
				UpdateTime:   now,
				CreationTime: now,
				Name:         "order 1",
				FullName:     "test test",
				Address:      "via roma 1",
				Goods: []model.GoodStock{
					{
						GoodID:   "1",
						Quantity: 1,
					},
				},
				Reservations: []string{"1"},
				Warehouses:   []model.OrderWarehouseUsed{},
			}, nil)
			suite.setCompletedWarehousePort.EXPECT().SetCompletedWarehouse(gomock.Any()).Return(model.Order{
				ID:           "1",
				Status:       "Filled",
				UpdateTime:   now,
				CreationTime: now,
				Name:         "order 1",
				FullName:     "test test",
				Address:      "via roma 1",
				Goods: []model.GoodStock{
					{
						GoodID:   "1",
						Quantity: 1,
					},
				},
				Reservations: []string{"1"},
				Warehouses:   []model.OrderWarehouseUsed{{WarehouseID: "1", Goods: map[string]int64{"1": 1}}},
			}, nil)
			suite.setCompletedWarehousePort.EXPECT().SetComplete(gomock.Any()).Return(nil)
		},
		func() fx.Option { return fx.Options() },
		func() interface{} {
			return func(service *ApplyStockUpdateService) {
				cmd := port.StockUpdateCmd{
					ID:            "1",
					Type:          port.StockUpdateCmdTypeOrder,
					WarehouseID:   "1",
					Timestamp:     time.Now().UnixMilli(),
					ReservationID: "1",
					TransferID:    "",
					OrderID:       "1",
					Goods: []port.StockUpdateGood{
						{
							GoodID:   "1",
							Quantity: 1,
						},
					},
				}
				err := service.ApplyStockUpdate(ctx, cmd)
				require.NoError(t, err)
			}
		},
	)
}

func TestApplyStockUpdateServiceTransfer(t *testing.T) {
	ctx := t.Context()
	runTestApplyStockUpdateService(t,
		func(suite *applyStockUpdateServiceMockSuite) {
			now := time.Now().UnixMilli()
			suite.applyStockUpdatePort.EXPECT().ApplyStockUpdate(gomock.Any())
			firstCall := suite.getTransferPort.EXPECT().GetTransfer(gomock.Any()).Return(model.Transfer{
				ID:                "1",
				Status:            "Filled",
				SenderID:          "1",
				ReceiverID:        "2",
				ReservationID:     "1",
				CreationTime:      now,
				UpdateTime:        now,
				LinkedStockUpdate: 1,
				Goods: []model.GoodStock{
					{
						GoodID:   "1",
						Quantity: 1,
					},
				},
			}, nil)

			suite.getTransferPort.EXPECT().GetTransfer(gomock.Any()).After(firstCall).Return(model.Transfer{
				ID:                "1",
				Status:            "Filled",
				SenderID:          "1",
				ReceiverID:        "2",
				ReservationID:     "1",
				CreationTime:      now,
				UpdateTime:        now,
				LinkedStockUpdate: 2,
				Goods: []model.GoodStock{
					{
						GoodID:   "1",
						Quantity: 1,
					},
				},
			}, nil)
			suite.setCompleteTransferPort.EXPECT().IncrementLinkedStockUpdate(gomock.Any()).Return(nil)
			suite.setCompleteTransferPort.EXPECT().SetComplete(gomock.Any()).Return(nil)
		},
		func() fx.Option { return fx.Options() },
		func() interface{} {
			return func(service *ApplyStockUpdateService) {
				cmd := port.StockUpdateCmd{
					ID:            "1",
					Type:          port.StockUpdateCmdTypeTransfer,
					WarehouseID:   "1",
					Timestamp:     time.Now().UnixMilli(),
					ReservationID: "1",
					TransferID:    "1",
					OrderID:       "",
					Goods: []port.StockUpdateGood{
						{
							GoodID:   "1",
							Quantity: 1,
						},
					},
				}
				err := service.ApplyStockUpdate(ctx, cmd)
				require.NoError(t, err)
			}
		},
	)
}
