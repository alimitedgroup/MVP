package business

import (
	"errors"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	gomock "go.uber.org/mock/gomock"
)

type managerOrderServiceMockSuite struct {
	getOrderPort                 *MockIGetOrderPort
	getTransferPort              *MockIGetTransferPort
	sendOrderUpdatePort          *MockISendOrderUpdatePort
	sendTransferUpdatePort       *MockISendTransferUpdatePort
	sendContactWarehousePort     *MockISendContactWarehousePort
	requestReservationPort       *MockIRequestReservationPort
	calculateAvailabilityUseCase *MockICalculateAvailabilityUseCase
}

func newManagerOrderServiceMockSuite(t *testing.T) *managerOrderServiceMockSuite {
	ctrl := gomock.NewController(t)

	return &managerOrderServiceMockSuite{
		getOrderPort:                 NewMockIGetOrderPort(ctrl),
		getTransferPort:              NewMockIGetTransferPort(ctrl),
		sendOrderUpdatePort:          NewMockISendOrderUpdatePort(ctrl),
		sendTransferUpdatePort:       NewMockISendTransferUpdatePort(ctrl),
		sendContactWarehousePort:     NewMockISendContactWarehousePort(ctrl),
		requestReservationPort:       NewMockIRequestReservationPort(ctrl),
		calculateAvailabilityUseCase: NewMockICalculateAvailabilityUseCase(ctrl),
	}
}

func runTestManagerOrderService(t *testing.T, build func(*managerOrderServiceMockSuite), buildOptions func() fx.Option, runLifeCycle func() interface{}) {
	ctx := t.Context()
	suite := newManagerOrderServiceMockSuite(t)

	build(suite)

	app := fx.New(
		fx.Supply(fx.Annotate(suite.getOrderPort, fx.As(new(port.IGetOrderPort)))),
		fx.Supply(fx.Annotate(suite.getTransferPort, fx.As(new(port.IGetTransferPort)))),
		fx.Supply(fx.Annotate(suite.sendOrderUpdatePort, fx.As(new(port.ISendOrderUpdatePort)))),
		fx.Supply(fx.Annotate(suite.sendTransferUpdatePort, fx.As(new(port.ISendTransferUpdatePort)))),
		fx.Supply(fx.Annotate(suite.sendContactWarehousePort, fx.As(new(port.ISendContactWarehousePort)))),
		fx.Supply(fx.Annotate(suite.requestReservationPort, fx.As(new(port.IRequestReservationPort)))),
		fx.Supply(fx.Annotate(suite.calculateAvailabilityUseCase, fx.As(new(port.ICalculateAvailabilityUseCase)))),
		fx.Provide(NewManageOrderService),
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

func TestManageOrderServiceGetAllTransfers(t *testing.T) {
	ctx := t.Context()
	runTestManagerOrderService(t,
		func(suite *managerOrderServiceMockSuite) {
			suite.getTransferPort.EXPECT().GetAllTransfer().Return([]model.Transfer{
				{
					ID:                "1",
					SenderID:          "1",
					ReceiverID:        "2",
					Status:            "Created",
					UpdateTime:        0,
					CreationTime:      0,
					LinkedStockUpdate: 0,
					ReservationID:     "",
					Goods: []model.GoodStock{
						{
							GoodID:   "1",
							Quantity: 1,
						},
					},
				},
			})
		},
		func() fx.Option { return fx.Options() },
		func() interface{} {
			return func(service *ManageOrderService) {
				transfers := service.GetAllTransfers(ctx)
				require.Len(t, transfers, 1)
			}
		},
	)

}

func TestManageOrderServiceGetTransfer(t *testing.T) {
	ctx := t.Context()
	runTestManagerOrderService(t,
		func(suite *managerOrderServiceMockSuite) {
			suite.getTransferPort.EXPECT().GetTransfer(gomock.Any()).Return(model.Transfer{

				ID:                "1",
				SenderID:          "1",
				ReceiverID:        "2",
				Status:            "Created",
				UpdateTime:        0,
				CreationTime:      0,
				LinkedStockUpdate: 0,
				ReservationID:     "",
				Goods: []model.GoodStock{
					{
						GoodID:   "1",
						Quantity: 1,
					},
				},
			}, nil)
		},
		func() fx.Option { return fx.Options() },
		func() interface{} {
			return func(service *ManageOrderService) {
				transfer, err := service.GetTransfer(ctx, "1")
				require.NoError(t, err)
				require.Equal(t, transfer.ID, "1")
			}
		},
	)

}

func TestManageOrderServiceGetAllOrders(t *testing.T) {
	ctx := t.Context()
	runTestManagerOrderService(t,
		func(suite *managerOrderServiceMockSuite) {
			suite.getOrderPort.EXPECT().GetAllOrder().Return([]model.Order{
				{
					ID:           "1",
					Name:         "order 1",
					FullName:     "test test",
					Address:      "via roma 1",
					Status:       "Created",
					UpdateTime:   0,
					CreationTime: 0,
					Reservations: []string{},
					Warehouses:   []model.OrderWarehouseUsed{},
					Goods: []model.GoodStock{
						{
							GoodID:   "1",
							Quantity: 1,
						},
					},
				},
			})
		},
		func() fx.Option { return fx.Options() },
		func() interface{} {
			return func(service *ManageOrderService) {
				orders := service.GetAllOrders(ctx)
				require.Len(t, orders, 1)
			}
		},
	)

}

func TestManageOrderServiceGetOrder(t *testing.T) {
	ctx := t.Context()
	runTestManagerOrderService(t,
		func(suite *managerOrderServiceMockSuite) {
			suite.getOrderPort.EXPECT().GetOrder(gomock.Any()).Return(model.Order{
				ID:           "1",
				Status:       "Created",
				UpdateTime:   0,
				CreationTime: 0,
				Name:         "order 1",
				FullName:     "test test",
				Address:      "via roma 1",
				Reservations: []string{},
				Warehouses:   []model.OrderWarehouseUsed{},
				Goods: []model.GoodStock{
					{
						GoodID:   "1",
						Quantity: 1,
					},
				},
			}, nil)
		},
		func() fx.Option { return fx.Options() },
		func() interface{} {
			return func(service *ManageOrderService) {
				order, err := service.GetOrder(ctx, "1")
				require.NoError(t, err)
				require.Equal(t, order.ID, "1")
			}
		},
	)

}

func TestManageOrderServiceCreateOrder(t *testing.T) {
	ctx := t.Context()
	runTestManagerOrderService(t,
		func(suite *managerOrderServiceMockSuite) {
			suite.sendOrderUpdatePort.EXPECT().SendOrderUpdate(gomock.Any(), gomock.Any()).Return(model.Order{
				ID:           "1",
				Status:       "Created",
				UpdateTime:   0,
				CreationTime: 0,
				Name:         "order 1",
				FullName:     "test test",
				Address:      "via roma 1",
				Reservations: []string{},
				Warehouses:   []model.OrderWarehouseUsed{},
				Goods: []model.GoodStock{
					{
						GoodID:   "1",
						Quantity: 1,
					},
				},
			}, nil)
			suite.sendContactWarehousePort.EXPECT().SendContactWarehouses(gomock.Any(), gomock.Any()).Return(nil)
		},
		func() fx.Option { return fx.Options() },
		func() interface{} {
			return func(service *ManageOrderService) {
				cmd := port.CreateOrderCmd{
					Name:     "order 1",
					FullName: "test test",
					Address:  "via roma 1",
					Goods: []port.CreateOrderGood{
						{
							GoodID:   "1",
							Quantity: 1,
						},
					},
				}
				resp, err := service.CreateOrder(ctx, cmd)
				require.NoError(t, err)
				require.NotEmpty(t, resp.OrderID)
			}
		},
	)
}

func TestManageOrderServiceCreateTrasfer(t *testing.T) {
	ctx := t.Context()
	runTestManagerOrderService(t,
		func(suite *managerOrderServiceMockSuite) {
			suite.sendTransferUpdatePort.EXPECT().SendTransferUpdate(gomock.Any(), gomock.Any()).Return(model.Transfer{
				ID:                "1",
				SenderID:          "1",
				ReceiverID:        "2",
				Status:            "Created",
				UpdateTime:        0,
				CreationTime:      0,
				LinkedStockUpdate: 0,
				ReservationID:     "",
				Goods: []model.GoodStock{
					{
						GoodID:   "1",
						Quantity: 1,
					},
				},
			}, nil)
			suite.sendContactWarehousePort.EXPECT().SendContactWarehouses(gomock.Any(), gomock.Any()).Return(nil)
		},
		func() fx.Option { return fx.Options() },
		func() interface{} {
			return func(service *ManageOrderService) {
				cmd := port.CreateTransferCmd{
					SenderID:   "1",
					ReceiverID: "2",
					Goods: []port.CreateTransferGood{
						{
							GoodID:   "1",
							Quantity: 1,
						},
					},
				}
				resp, err := service.CreateTransfer(ctx, cmd)
				require.NoError(t, err)
				require.NotEmpty(t, resp.TransferID)
			}
		},
	)
}

func TestManageOrderServiceContactWarehouseTransfer(t *testing.T) {
	ctx := t.Context()
	runTestManagerOrderService(t,
		func(suite *managerOrderServiceMockSuite) {
			suite.sendTransferUpdatePort.EXPECT().SendTransferUpdate(gomock.Any(), gomock.Any()).Return(model.Transfer{
				ID:                "1",
				SenderID:          "1",
				ReceiverID:        "2",
				Status:            "Filled",
				UpdateTime:        0,
				CreationTime:      0,
				LinkedStockUpdate: 1,
				ReservationID:     "1",
				Goods: []model.GoodStock{
					{
						GoodID:   "1",
						Quantity: 1,
					},
				},
			}, nil)
			suite.requestReservationPort.EXPECT().RequestReservation(gomock.Any(), gomock.Any()).Return(port.RequestReservationResponse{ID: "1"}, nil)
		},
		func() fx.Option { return fx.Options() },
		func() interface{} {
			return func(service *ManageOrderService) {
				cmd := port.ContactWarehousesCmd{
					Type:  port.ContactWarehousesTypeTransfer,
					Order: nil,
					Transfer: &port.ContactWarehousesTransfer{
						ID:         "1",
						SenderID:   "1",
						ReceiverID: "2",
						Status:     "Created",
						Goods: []port.ContactWarehousesGood{
							{
								GoodID:   "1",
								Quantity: 1,
							},
						},
						UpdateTime:    0,
						CreationTime:  0,
						ReservationID: "",
					},
					RetryUntil:            0,
					RetryInTime:           0,
					ExcludeWarehouses:     []string{},
					ConfirmedReservations: []port.ConfirmedReservation{},
				}
				resp, err := service.ContactWarehouses(ctx, cmd)
				require.NoError(t, err)
				require.False(t, resp.IsRetry)
			}
		},
	)
}

func TestManageOrderServiceContactWarehouseTransferRetryLater(t *testing.T) {
	ctx := t.Context()
	runTestManagerOrderService(t,
		func(suite *managerOrderServiceMockSuite) {
			suite.requestReservationPort.EXPECT().RequestReservation(gomock.Any(), gomock.Any()).Return(port.RequestReservationResponse{}, errors.New("connection error"))
			suite.sendContactWarehousePort.EXPECT().SendContactWarehouses(gomock.Any(), gomock.Any()).Return(nil)
		},
		func() fx.Option { return fx.Options() },
		func() interface{} {
			return func(service *ManageOrderService) {
				cmd := port.ContactWarehousesCmd{
					Type:  port.ContactWarehousesTypeTransfer,
					Order: nil,
					Transfer: &port.ContactWarehousesTransfer{
						ID:         "1",
						SenderID:   "1",
						ReceiverID: "2",
						Status:     "Created",
						Goods: []port.ContactWarehousesGood{
							{
								GoodID:   "1",
								Quantity: 1,
							},
						},
						UpdateTime:   0,
						CreationTime: 0,
					},
					RetryUntil:            time.Now().Add(time.Hour).UnixMilli(),
					RetryInTime:           (time.Minute * 2).Milliseconds(),
					ExcludeWarehouses:     []string{},
					ConfirmedReservations: []port.ConfirmedReservation{},
				}
				resp, err := service.ContactWarehouses(ctx, cmd)
				require.NoError(t, err)
				require.False(t, resp.IsRetry)
			}
		},
	)
}

func TestManageOrderServiceContactWarehouseOrder(t *testing.T) {
	ctx := t.Context()
	runTestManagerOrderService(t,
		func(suite *managerOrderServiceMockSuite) {
			suite.calculateAvailabilityUseCase.EXPECT().GetAvailable(gomock.Any(), gomock.Any()).Return(
				port.CalculateAvailabilityResponse{
					Warehouses: []port.WarehouseAvailability{
						{
							WarehouseID: "1",
							Goods:       map[string]int64{"1": 1},
						},
					},
				}, nil)
			suite.requestReservationPort.EXPECT().RequestReservation(gomock.Any(), gomock.Any()).Return(port.RequestReservationResponse{ID: "1"}, nil)
			suite.sendOrderUpdatePort.EXPECT().SendOrderUpdate(gomock.Any(), gomock.Any()).Return(model.Order{
				ID:           "1",
				Status:       "Filled",
				Name:         "order 1",
				FullName:     "test test",
				Address:      "via roma 1",
				Reservations: []string{"1"},
				Warehouses:   []model.OrderWarehouseUsed{},
				UpdateTime:   0,
				CreationTime: 0,
				Goods: []model.GoodStock{
					{
						GoodID:   "1",
						Quantity: 1,
					},
				},
			}, nil)
		},
		func() fx.Option { return fx.Options() },
		func() interface{} {
			return func(service *ManageOrderService) {
				cmd := port.ContactWarehousesCmd{
					Type: port.ContactWarehousesTypeOrder,
					Order: &port.ContactWarehousesOrder{
						ID:           "1",
						Status:       "Created",
						Name:         "order 1",
						FullName:     "test test",
						Address:      "via roma 1",
						Reservations: []string{},
						Goods: []port.ContactWarehousesGood{
							{
								GoodID:   "1",
								Quantity: 1,
							},
						},
						UpdateTime:   0,
						CreationTime: 0,
					},
					Transfer:              nil,
					RetryUntil:            0,
					RetryInTime:           0,
					ExcludeWarehouses:     []string{},
					ConfirmedReservations: []port.ConfirmedReservation{},
				}
				resp, err := service.ContactWarehouses(ctx, cmd)
				require.NoError(t, err)
				require.False(t, resp.IsRetry)
			}
		},
	)
}

func TestManageOrderServiceContactWarehouseOrderRetryLater(t *testing.T) {
	ctx := t.Context()
	runTestManagerOrderService(t,
		func(suite *managerOrderServiceMockSuite) {
			suite.calculateAvailabilityUseCase.EXPECT().GetAvailable(gomock.Any(), gomock.Any()).Return(
				port.CalculateAvailabilityResponse{
					Warehouses: []port.WarehouseAvailability{
						{
							WarehouseID: "1",
							Goods:       map[string]int64{"1": 1},
						},
					},
				}, nil)
			suite.requestReservationPort.EXPECT().RequestReservation(gomock.Any(), gomock.Any()).Return(port.RequestReservationResponse{}, errors.New("connectio error"))
			suite.sendContactWarehousePort.EXPECT().SendContactWarehouses(gomock.Any(), gomock.Any()).Return(nil)
		},
		func() fx.Option { return fx.Options() },
		func() interface{} {
			return func(service *ManageOrderService) {
				cmd := port.ContactWarehousesCmd{
					Type: port.ContactWarehousesTypeOrder,
					Order: &port.ContactWarehousesOrder{
						ID:           "1",
						Status:       "Created",
						Name:         "order 1",
						FullName:     "test test",
						Address:      "via roma 1",
						Reservations: []string{},
						Goods: []port.ContactWarehousesGood{
							{
								GoodID:   "1",
								Quantity: 1,
							},
						},
						UpdateTime:   0,
						CreationTime: 0,
					},
					Transfer:              nil,
					RetryUntil:            time.Now().Add(time.Hour).UnixMilli(),
					RetryInTime:           (time.Minute * 2).Milliseconds(),
					ExcludeWarehouses:     []string{},
					ConfirmedReservations: []port.ConfirmedReservation{},
				}
				resp, err := service.ContactWarehouses(ctx, cmd)
				require.NoError(t, err)
				require.False(t, resp.IsRetry)
			}
		},
	)
}
