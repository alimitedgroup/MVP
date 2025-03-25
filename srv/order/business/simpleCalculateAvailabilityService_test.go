package business

import (
	"testing"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	gomock "go.uber.org/mock/gomock"
)

type simpleCalculateAvailabilityServiceMockSuite struct {
	getStockPort *MockIGetStockPort
}

func newSimpleCalculateAvailabilityServiceMockSuite(t *testing.T) *simpleCalculateAvailabilityServiceMockSuite {
	ctrl := gomock.NewController(t)

	return &simpleCalculateAvailabilityServiceMockSuite{
		getStockPort: NewMockIGetStockPort(ctrl),
	}
}

func runTestsimpleCalculateAvailabilityService(t *testing.T, build func(*simpleCalculateAvailabilityServiceMockSuite), buildOptions func() fx.Option, runLifeCycle func() interface{}) {
	ctx := t.Context()
	suite := newSimpleCalculateAvailabilityServiceMockSuite(t)

	build(suite)

	app := fx.New(
		fx.Supply(fx.Annotate(suite.getStockPort, fx.As(new(port.IGetStockPort)))),
		fx.Provide(NewSimpleCalculateAvailabilityService),
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

func TestSimpleCalculateAvailabilityService(t *testing.T) {
	ctx := t.Context()
	runTestsimpleCalculateAvailabilityService(t,
		func(suite *simpleCalculateAvailabilityServiceMockSuite) {
			suite.getStockPort.EXPECT().GetGlobalStock(gomock.Any()).Return(model.GoodStock{GoodID: "1", Quantity: 1})
			suite.getStockPort.EXPECT().GetWarehouses().Return([]model.Warehouse{{ID: "1"}})
			suite.getStockPort.EXPECT().GetStock(gomock.Any()).Return(model.GoodStock{GoodID: "1", Quantity: 1}, nil)
		},
		func() fx.Option { return fx.Options() },
		func() interface{} {
			return func(srv *SimpleCalculateAvailabilityService) {
				cmd := port.CalculateAvailabilityCmd{
					Goods: []port.RequestedGood{
						{GoodID: "1", Quantity: 1},
					},
					ExcludedWarehouses: []string{},
				}
				resp, err := srv.GetAvailable(ctx, cmd)
				require.NoError(t, err)
				require.Equal(t, resp.Warehouses[0].WarehouseID, "1")
				require.Equal(t, resp.Warehouses[0].Goods, map[string]int64{"1": 1})
			}
		},
	)
}
