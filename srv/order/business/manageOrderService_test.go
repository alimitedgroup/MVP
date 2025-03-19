package business

import (
	"testing"

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
					Id:                "1",
					SenderId:          "1",
					ReceiverId:        "2",
					Status:            "Created",
					UpdateTime:        0,
					CreationTime:      0,
					LinkedStockUpdate: 0,
					ReservationID:     "",
					Goods: []model.GoodStock{
						{
							ID:       "1",
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
