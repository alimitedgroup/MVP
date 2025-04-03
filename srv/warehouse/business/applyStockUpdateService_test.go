package business

import (
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
)

func TestApplyStockUpdateService(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)

	idempotentPortMock := NewMockIIdempotentPort(ctrl)
	idempotentPortMock.EXPECT().SaveEventID(gomock.Any())

	applyStockUpdatePortmock := NewMockIApplyStockUpdatePort(ctrl)
	applyStockUpdatePortmock.EXPECT().ApplyStockUpdate(gomock.Any())

	transactionPort := NewMockITransactionPort(ctrl)
	transactionPort.EXPECT().Lock()
	transactionPort.EXPECT().Unlock()

	app := fx.New(
		fx.Supply(fx.Annotate(idempotentPortMock, fx.As(new(port.IIdempotentPort)))),
		fx.Supply(fx.Annotate(applyStockUpdatePortmock, fx.As(new(port.IApplyStockUpdatePort)))),
		fx.Supply(fx.Annotate(transactionPort, fx.As(new(port.ITransactionPort)))),
		fx.Provide(NewApplyStockUpdateService),
		fx.Invoke(func(service *ApplyStockUpdateService) {
			cmd := port.StockUpdateCmd{
				Type:          port.StockUpdateCmdTypeOrder,
				OrderID:       "1",
				TransferID:    "",
				ReservationID: "1",
				Timestamp:     time.Now().UnixMilli(),
				ID:            "1",
				Goods: []port.StockUpdateGood{
					{
						GoodID:   "1",
						Quantity: 10,
					},
					{
						GoodID:   "2",
						Quantity: 20,
					},
				},
			}

			service.ApplyStockUpdate(cmd)
		}),
	)

	err := app.Start(ctx)
	require.NoError(t, err)

	defer func() {
		err := app.Stop(ctx)
		require.NoError(t, err)
	}()

}
