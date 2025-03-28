package listener

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/lib"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/common/stream"
	internalStream "github.com/alimitedgroup/MVP/srv/order/adapter/stream"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	gomock "go.uber.org/mock/gomock"
)

type orderListenerMockSuite struct {
	applyOrderUpdateUseCaseMock    *MockIApplyOrderUpdateUseCase
	applyTransferUpdateUseCaseMock *MockIApplyTransferUpdateUseCase
	contactWarehouseUseCaseMock    *MockIContactWarehousesUseCase
}

func NewOrderListenerMockSuite(t *testing.T) *orderListenerMockSuite {
	ctrl := gomock.NewController(t)

	return &orderListenerMockSuite{
		applyOrderUpdateUseCaseMock:    NewMockIApplyOrderUpdateUseCase(ctrl),
		applyTransferUpdateUseCaseMock: NewMockIApplyTransferUpdateUseCase(ctrl),
		contactWarehouseUseCaseMock:    NewMockIContactWarehousesUseCase(ctrl),
	}
}

func runTestOrderListener(t *testing.T, build func(*orderListenerMockSuite), buildOptions func() fx.Option, runLifeCycle func() interface{}) {
	ctx := t.Context()
	suite := NewOrderListenerMockSuite(t)

	build(suite)

	app := fx.New(
		fx.Supply(fx.Annotate(suite.applyOrderUpdateUseCaseMock, fx.As(new(port.IApplyOrderUpdateUseCase)))),
		fx.Supply(fx.Annotate(suite.applyTransferUpdateUseCaseMock, fx.As(new(port.IApplyTransferUpdateUseCase)))),
		fx.Supply(fx.Annotate(suite.contactWarehouseUseCaseMock, fx.As(new(port.IContactWarehousesUseCase)))),
		fx.Provide(observability.TestMeter),
		fx.Provide(observability.TestLogger),
		fx.Provide(NewOrderListener),
		fx.Provide(NewOrderRouter),
		fx.Invoke(runLifeCycle()),
		buildOptions(),
		lib.ModuleTest,
		fx.Supply(t),
	)

	err := app.Start(ctx)
	require.NoError(t, err)

	defer func() {
		err := app.Stop(ctx)
		require.NoError(t, err)
	}()
}

func TestOrderListenerApplyOrderUpdate(t *testing.T) {
	ns, _ := broker.NewInProcessNATSServer(t)
	js, err := jetstream.New(ns)
	require.NoError(t, err)

	runTestOrderListener(t,
		func(suite *orderListenerMockSuite) {
			suite.applyOrderUpdateUseCaseMock.EXPECT().ApplyOrderUpdate(gomock.Any(), gomock.Any())
		},
		func() fx.Option {
			return fx.Options(fx.Supply(ns))
		},
		func() interface{} {
			return func(lc fx.Lifecycle, r *OrderRouter) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						err := r.Setup(ctx)
						require.NoError(t, err)

						event := stream.OrderUpdate{
							ID:           "1",
							Status:       "Created",
							Name:         "Test",
							FullName:     "Test Test",
							Address:      "via roma 11",
							Reservations: []string{},
							CreationTime: time.Now().UnixMilli(),
							UpdateTime:   time.Now().UnixMilli(),
							Goods:        []stream.OrderUpdateGood{{GoodID: "1", Quantity: 1}},
						}
						payload, err := json.Marshal(event)
						require.NoError(t, err)

						resp, err := js.Publish(ctx, "order.update", payload)
						require.NoError(t, err)
						require.Equal(t, resp.Stream, "order_update")

						time.Sleep(100 * time.Millisecond)

						return nil
					},
				})
			}
		},
	)
}

func TestOrderListenerApplyTransferUpdate(t *testing.T) {
	ns, _ := broker.NewInProcessNATSServer(t)
	js, err := jetstream.New(ns)
	require.NoError(t, err)

	runTestOrderListener(t,
		func(suite *orderListenerMockSuite) {
			suite.applyTransferUpdateUseCaseMock.EXPECT().ApplyTransferUpdate(gomock.Any(), gomock.Any())
		},
		func() fx.Option {
			return fx.Options(fx.Supply(ns))
		},
		func() interface{} {
			return func(lc fx.Lifecycle, r *OrderRouter) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						err := r.Setup(ctx)
						require.NoError(t, err)

						event := stream.TransferUpdate{
							ID:            "1",
							Status:        "Created",
							SenderID:      "1",
							ReceiverID:    "2",
							ReservationID: "",
							CreationTime:  time.Now().UnixMilli(),
							UpdateTime:    time.Now().UnixMilli(),
							Goods:         []stream.TransferUpdateGood{{GoodID: "1", Quantity: 1}},
						}
						payload, err := json.Marshal(event)
						require.NoError(t, err)

						resp, err := js.Publish(ctx, "transfer.update", payload)
						require.NoError(t, err)
						require.Equal(t, resp.Stream, "transfer_update")

						time.Sleep(100 * time.Millisecond)

						return nil
					},
				})
			}
		},
	)
}

func TestOrderListenerContactWarehousesTransfer(t *testing.T) {
	ns, _ := broker.NewInProcessNATSServer(t)
	js, err := jetstream.New(ns)
	require.NoError(t, err)

	runTestOrderListener(t,
		func(suite *orderListenerMockSuite) {
			suite.contactWarehouseUseCaseMock.EXPECT().ContactWarehouses(gomock.Any(), gomock.Any()).Return(port.ContactWarehousesResponse{IsRetry: false}, nil)
		},
		func() fx.Option {
			return fx.Options(fx.Supply(ns))
		},
		func() interface{} {
			return func(lc fx.Lifecycle, r *OrderRouter) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						err := r.Setup(ctx)
						require.NoError(t, err)

						event := internalStream.ContactWarehouses{
							Order: nil,
							Transfer: &internalStream.ContactWarehousesTransfer{
								ID:            "1",
								Status:        "Created",
								SenderID:      "1",
								ReceiverID:    "2",
								ReservationID: "",
								CreationTime:  time.Now().UnixMilli(),
								UpdateTime:    time.Now().UnixMilli(),
								Goods:         []internalStream.ContactWarehousesGood{{GoodID: "1", Quantity: 1}},
							},
							Type:              internalStream.ContactWarehousesTypeTransfer,
							ExcludeWarehouses: []string{},
							RetryInTime:       0,
							RetryUntil:        time.Now().Add(1 * time.Hour).UnixMilli(),
						}
						payload, err := json.Marshal(event)
						require.NoError(t, err)

						resp, err := js.Publish(ctx, "contact.warehouses", payload)
						require.NoError(t, err)
						require.Equal(t, resp.Stream, "contact_warehouses")
						time.Sleep(100 * time.Millisecond)
						return nil
					},
				})
			}
		},
	)
}

func TestOrderListenerContactWarehousesOrder(t *testing.T) {
	ns, _ := broker.NewInProcessNATSServer(t)
	js, err := jetstream.New(ns)
	require.NoError(t, err)

	runTestOrderListener(t,
		func(suite *orderListenerMockSuite) {
			suite.contactWarehouseUseCaseMock.EXPECT().ContactWarehouses(gomock.Any(), gomock.Any()).Return(port.ContactWarehousesResponse{IsRetry: false}, nil)
		},
		func() fx.Option {
			return fx.Options(fx.Supply(ns))
		},
		func() interface{} {
			return func(lc fx.Lifecycle, r *OrderRouter) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						err := r.Setup(ctx)
						require.NoError(t, err)

						event := internalStream.ContactWarehouses{
							Transfer: nil,
							Order: &internalStream.ContactWarehousesOrder{
								ID:           "1",
								Status:       "Created",
								Name:         "Test",
								FullName:     "Test Test",
								Address:      "via roma 11",
								Reservations: []string{},
								CreationTime: time.Now().UnixMilli(),
								UpdateTime:   time.Now().UnixMilli(),
								Goods:        []internalStream.ContactWarehousesGood{{GoodID: "1", Quantity: 1}},
							},
							Type:              internalStream.ContactWarehousesTypeOrder,
							ExcludeWarehouses: []string{},
							RetryInTime:       0,
							RetryUntil:        time.Now().Add(1 * time.Hour).UnixMilli(),
						}
						payload, err := json.Marshal(event)
						require.NoError(t, err)

						resp, err := js.Publish(ctx, "contact.warehouses", payload)
						require.NoError(t, err)
						require.Equal(t, resp.Stream, "contact_warehouses")

						time.Sleep(100 * time.Millisecond)

						return nil
					},
				})
			}
		},
	)
}

func TestOrderListenerContactWarehousesOrderWithRetry(t *testing.T) {
	ns, _ := broker.NewInProcessNATSServer(t)
	js, err := jetstream.New(ns)
	require.NoError(t, err)

	runTestOrderListener(t,
		func(suite *orderListenerMockSuite) {
			suite.contactWarehouseUseCaseMock.EXPECT().ContactWarehouses(gomock.Any(), gomock.Any()).Return(port.ContactWarehousesResponse{IsRetry: true, RetryAfter: 1 * time.Hour}, nil)
		},
		func() fx.Option {
			return fx.Options(fx.Supply(ns))
		},
		func() interface{} {
			return func(lc fx.Lifecycle, r *OrderRouter) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						err := r.Setup(ctx)
						require.NoError(t, err)

						event := internalStream.ContactWarehouses{
							Transfer: nil,
							Order: &internalStream.ContactWarehousesOrder{
								ID:           "1",
								Status:       "Created",
								Name:         "Test",
								FullName:     "Test Test",
								Address:      "via roma 11",
								Reservations: []string{},
								CreationTime: time.Now().UnixMilli(),
								UpdateTime:   time.Now().UnixMilli(),
								Goods:        []internalStream.ContactWarehousesGood{{GoodID: "1", Quantity: 1}},
							},
							Type:              internalStream.ContactWarehousesTypeOrder,
							ExcludeWarehouses: []string{},
							RetryInTime:       (1 * time.Hour).Milliseconds(),
							RetryUntil:        time.Now().Add(1 * time.Hour).UnixMilli(),
						}
						payload, err := json.Marshal(event)
						require.NoError(t, err)

						resp, err := js.Publish(ctx, "contact.warehouses", payload)
						require.NoError(t, err)
						require.Equal(t, resp.Stream, "contact_warehouses")

						time.Sleep(100 * time.Millisecond)

						return nil
					},
				})
			}
		},
	)
}
