package adapterout

import (
	"testing"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestGetTransfers(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)

	sub, err := nc.Subscribe("transfer.get.all", func(msg *nats.Msg) {
		err := msg.Respond([]byte(`{"error": "", "message": []}`))
		require.NoError(t, err)
	})
	require.NoError(t, err)
	defer func() {
		err := sub.Unsubscribe()
		require.NoError(t, err)
	}()

	brk, err := broker.NewNatsMessageBroker(nc, zaptest.NewLogger(t))
	require.NoError(t, err)
	order := NewOrderAdapter(brk)

	info, err := order.GetAllTransfers()
	require.NoError(t, err)
	require.Len(t, info, 0)
}

func TestGetOrders(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)

	sub, err := nc.Subscribe("order.get.all", func(msg *nats.Msg) {
		err := msg.Respond([]byte(`{"error": "", "message": []}`))
		require.NoError(t, err)
	})
	require.NoError(t, err)
	defer func() {
		err := sub.Unsubscribe()
		require.NoError(t, err)
	}()

	brk, err := broker.NewNatsMessageBroker(nc, zaptest.NewLogger(t))
	require.NoError(t, err)
	order := NewOrderAdapter(brk)

	info, err := order.GetAllOrders()
	require.NoError(t, err)
	require.Len(t, info, 0)
}

func TestCreateOrder(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)

	sub, err := nc.Subscribe("order.create", func(msg *nats.Msg) {
		err := msg.Respond([]byte(`{"error": "", "message": {"order_id": "1"}}`))
		require.NoError(t, err)
	})
	require.NoError(t, err)
	defer func() {
		err := sub.Unsubscribe()
		require.NoError(t, err)
	}()

	brk, err := broker.NewNatsMessageBroker(nc, zaptest.NewLogger(t))
	require.NoError(t, err)
	order := NewOrderAdapter(brk)

	dto := request.CreateOrderRequestDTO{
		Name:     "Order 1",
		FullName: "Mario Rossi",
		Address:  "Via Roma 1",
		Goods: []request.CreateOrderGood{
			{
				GoodID:   "1",
				Quantity: 1,
			},
		},
	}
	info, err := order.CreateOrder(dto)
	require.NoError(t, err)
	require.Equal(t, "1", info.OrderID)
}

func TestCreateTransfer(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)

	sub, err := nc.Subscribe("transfer.create", func(msg *nats.Msg) {
		err := msg.Respond([]byte(`{"error": "", "message": {"transfer_id": "1"}}`))
		require.NoError(t, err)
	})
	require.NoError(t, err)
	defer func() {
		err := sub.Unsubscribe()
		require.NoError(t, err)
	}()

	brk, err := broker.NewNatsMessageBroker(nc, zaptest.NewLogger(t))
	require.NoError(t, err)
	order := NewOrderAdapter(brk)

	dto := request.CreateTransferRequestDTO{
		SenderID:   "1",
		ReceiverID: "1",
		Goods: []request.TransferGood{
			{
				GoodID:   "1",
				Quantity: 1,
			},
		},
	}
	info, err := order.CreateTransfer(dto)
	require.NoError(t, err)
	require.Equal(t, "1", info.TransferID)
}
