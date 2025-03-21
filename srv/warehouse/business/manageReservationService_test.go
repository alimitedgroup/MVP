package business

import (
	"context"
	"testing"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
)

type manageReservationServiceMockSuite struct {
	createStockUpdatePortMock *MockICreateStockUpdatePort
	getStockPortMock          *MockIGetStockPort
	storeReservationEventPort *MockIStoreReservationEventPort
	getReservationPort        *MockIGetReservationPort
	applyReservationEventPort *MockIApplyReservationEventPort
	idempotentPortMock        *MockIIdempotentPort
}

func newManageReservationServiceMockSuite(t *testing.T) *manageReservationServiceMockSuite {
	ctrl := gomock.NewController(t)

	return &manageReservationServiceMockSuite{
		createStockUpdatePortMock: NewMockICreateStockUpdatePort(ctrl),
		getStockPortMock:          NewMockIGetStockPort(ctrl),
		storeReservationEventPort: NewMockIStoreReservationEventPort(ctrl),
		getReservationPort:        NewMockIGetReservationPort(ctrl),
		applyReservationEventPort: NewMockIApplyReservationEventPort(ctrl),
		idempotentPortMock:        NewMockIIdempotentPort(ctrl),
	}
}

func runTestManageReservationService(t *testing.T, build func(*manageReservationServiceMockSuite), buildOptions func() fx.Option, runLifeCycle func() interface{}) {
	ctx := t.Context()
	suite := newManageReservationServiceMockSuite(t)

	cfg := config.WarehouseConfig{ID: "1"}

	build(suite)

	app := fx.New(
		fx.Supply(&cfg),
		fx.Supply(fx.Annotate(suite.createStockUpdatePortMock, fx.As(new(port.ICreateStockUpdatePort)))),
		fx.Supply(fx.Annotate(suite.getStockPortMock, fx.As(new(port.IGetStockPort)))),
		fx.Supply(fx.Annotate(suite.storeReservationEventPort, fx.As(new(port.IStoreReservationEventPort)))),
		fx.Supply(fx.Annotate(suite.idempotentPortMock, fx.As(new(port.IIdempotentPort)))),
		fx.Supply(fx.Annotate(suite.applyReservationEventPort, fx.As(new(port.IApplyReservationEventPort)))),
		fx.Supply(fx.Annotate(suite.getReservationPort, fx.As(new(port.IGetReservationPort)))),
		fx.Provide(NewManageReservationService),
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

func TestManageReservationServiceApplyReservationEvent(t *testing.T) {
	runTestManageReservationService(t,
		func(suite *manageReservationServiceMockSuite) {
			suite.applyReservationEventPort.EXPECT().ApplyReservationEvent(gomock.Any()).Return(nil)
			suite.idempotentPortMock.EXPECT().IsAlreadyProcessed(gomock.Any()).Return(false)
		},
		func() fx.Option { return fx.Options() },
		func() interface{} {
			return func(lc fx.Lifecycle, service *ManageReservationService) {
				lc.Append(fx.Hook{OnStart: func(ctx context.Context) error {
					applyCmd := port.ApplyReservationEventCmd{
						Id: "1",
						Goods: []port.ReservationGood{
							{
								GoodID:   "1",
								Quantity: 10,
							},
						},
					}
					err := service.ApplyReservationEvent(applyCmd)
					require.NoError(t, err)

					return nil
				}})
			}
		})
}

func TestManageReservationServiceCreateReservation(t *testing.T) {
	runTestManageReservationService(t,
		func(suite *manageReservationServiceMockSuite) {
			suite.getStockPortMock.EXPECT().GetFreeStock(gomock.Any()).Return(model.GoodStock{
				ID:       "1",
				Quantity: 10,
			})
			suite.storeReservationEventPort.EXPECT().StoreReservationEvent(gomock.Any(), gomock.Any()).Return(nil)
		},
		func() fx.Option { return fx.Options() },
		func() interface{} {
			return func(lc fx.Lifecycle, service *ManageReservationService) {
				lc.Append(fx.Hook{OnStart: func(ctx context.Context) error {
					applyCmd := port.CreateReservationCmd{
						Goods: []port.ReservationGood{
							{
								GoodID:   "1",
								Quantity: 10,
							},
						},
					}
					resp, err := service.CreateReservation(ctx, applyCmd)
					require.NotEmpty(t, resp.ReservationID)
					require.NoError(t, err)

					return nil
				}})
			}
		})
}

func TestManageReservationServiceConfirmOrder(t *testing.T) {
	runTestManageReservationService(t,
		func(suite *manageReservationServiceMockSuite) {
			suite.getStockPortMock.EXPECT().GetStock(gomock.Any()).Return(model.GoodStock{
				ID:       "1",
				Quantity: 10,
			})
			suite.getReservationPort.EXPECT().GetReservation(gomock.Any()).Return(model.Reservation{
				ID: "1",
				Goods: []model.ReservationGood{
					{
						GoodID:   "1",
						Quantity: 10,
					},
				},
			}, nil)
			suite.createStockUpdatePortMock.EXPECT().CreateStockUpdate(gomock.Any(), gomock.Any()).Return(nil)
			suite.applyReservationEventPort.EXPECT().ApplyOrderFilled(gomock.Any()).Return(nil)
		},
		func() fx.Option { return fx.Options() },
		func() interface{} {
			return func(lc fx.Lifecycle, service *ManageReservationService) {
				lc.Append(fx.Hook{OnStart: func(ctx context.Context) error {
					confirmCmd := port.ConfirmOrderCmd{
						OrderID:      "1",
						Status:       "Filled",
						Reservations: []string{"1"},
						Goods: []port.OrderUpdateGood{
							{
								GoodID:   "1",
								Quantity: 10,
							},
						},
					}
					err := service.ConfirmOrder(ctx, confirmCmd)
					require.NoError(t, err)

					return nil
				}})
			}
		})
}

func TestManageReservationServiceConfirmTransferSender(t *testing.T) {
	runTestManageReservationService(t,
		func(suite *manageReservationServiceMockSuite) {
			suite.getStockPortMock.EXPECT().GetStock(gomock.Any()).Return(model.GoodStock{
				ID:       "1",
				Quantity: 10,
			})
			suite.getReservationPort.EXPECT().GetReservation(gomock.Any()).Return(model.Reservation{
				ID: "1",
				Goods: []model.ReservationGood{
					{
						GoodID:   "1",
						Quantity: 10,
					},
				},
			}, nil)
			suite.createStockUpdatePortMock.EXPECT().CreateStockUpdate(gomock.Any(), gomock.Any()).Return(nil)
			suite.applyReservationEventPort.EXPECT().ApplyOrderFilled(gomock.Any()).Return(nil)
		},
		func() fx.Option {
			return fx.Options(
				fx.Replace(&config.WarehouseConfig{ID: "1"}),
			)
		},
		func() interface{} {
			return func(lc fx.Lifecycle, service *ManageReservationService) {
				lc.Append(fx.Hook{OnStart: func(ctx context.Context) error {
					confirmCmd := port.ConfirmTransferCmd{
						TransferID:    "1",
						Status:        "Filled",
						SenderID:      "1",
						ReceiverID:    "2",
						ReservationId: "1",
						Goods: []port.TransferUpdateGood{
							{
								GoodID:   "1",
								Quantity: 10,
							},
						},
					}
					err := service.ConfirmTransfer(ctx, confirmCmd)
					require.NoError(t, err)

					return nil
				}})
			}
		})
}

func TestManageReservationServiceConfirmTransferReceiver(t *testing.T) {
	runTestManageReservationService(t,
		func(suite *manageReservationServiceMockSuite) {
			suite.getStockPortMock.EXPECT().GetStock(gomock.Any()).Return(model.GoodStock{
				ID:       "1",
				Quantity: 10,
			})
			suite.createStockUpdatePortMock.EXPECT().CreateStockUpdate(gomock.Any(), gomock.Any()).Return(nil)
		},
		func() fx.Option {
			return fx.Options(
				fx.Replace(&config.WarehouseConfig{ID: "2"}),
			)
		},
		func() interface{} {
			return func(lc fx.Lifecycle, service *ManageReservationService) {
				lc.Append(fx.Hook{OnStart: func(ctx context.Context) error {
					confirmCmd := port.ConfirmTransferCmd{
						TransferID:    "1",
						Status:        "Filled",
						SenderID:      "1",
						ReceiverID:    "2",
						ReservationId: "1",
						Goods: []port.TransferUpdateGood{
							{
								GoodID:   "1",
								Quantity: 10,
							},
						},
					}
					err := service.ConfirmTransfer(ctx, confirmCmd)
					require.NoError(t, err)

					return nil
				}})
			}
		})
}
