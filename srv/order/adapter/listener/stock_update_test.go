package listener

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/lib"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
)

func TestStockUpdateListener(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)

	applyStockUpdateUseCaseMock := NewMockIApplyStockUpdateUseCase(ctrl)
	applyStockUpdateUseCaseMock.EXPECT().ApplyStockUpdate(gomock.Any(), gomock.Any()).Return(nil)

	ns, _ := broker.NewInProcessNATSServer(t)
	js, err := jetstream.New(ns)
	require.NoError(t, err)

	app := fx.New(
		lib.ModuleTest,
		fx.Supply(ns, t),
		fx.Supply(fx.Annotate(applyStockUpdateUseCaseMock, fx.As(new(port.IApplyStockUpdateUseCase)))),
		fx.Provide(NewStockUpdateListener),
		fx.Provide(NewStockUpdateRouter),
		fx.Invoke(func(lc fx.Lifecycle, r *StockUpdateRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					require.NoError(t, err)

					event := stream.StockUpdate{
						ID:            "1",
						WarehouseID:   "1",
						Type:          stream.StockUpdateTypeAdd,
						OrderID:       "",
						TransferID:    "",
						ReservationID: "",
						Timestamp:     time.Now().UnixMilli(),
						Goods: []stream.StockUpdateGood{
							{
								GoodID:   "1",
								Quantity: 1,
							},
						},
					}
					payload, err := json.Marshal(event)
					require.NoError(t, err)

					resp, err := js.Publish(ctx, "stock.update.1", payload)
					require.NoError(t, err)
					require.Equal(t, resp.Stream, "stock_update")

					return nil
				},
			})
		}),
	)

	err = app.Start(ctx)
	require.NoError(t, err)

	defer func() {
		err = app.Stop(ctx)
		require.NoError(t, err)
	}()
}
