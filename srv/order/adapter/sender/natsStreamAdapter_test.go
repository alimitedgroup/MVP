package sender

import (
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	internalStream "github.com/alimitedgroup/MVP/srv/order/adapter/stream"
	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
)

func TestNatsStreamAdapterSendOrderUpdate(t *testing.T) {
	ctx := t.Context()

	ns, _ := broker.NewInProcessNATSServer(t)
	js, err := jetstream.New(ns)
	require.NoError(t, err)

	s, err := js.CreateStream(ctx, stream.OrderUpdateStreamConfig)
	require.NoError(t, err)

	broker, err := broker.NewNatsMessageBroker(ns)
	require.NoError(t, err)
	a := NewNatsStreamAdapter(broker)

	cmd := port.SendOrderUpdateCmd{
		ID: "1",
		Goods: []port.SendOrderUpdateGood{
			{
				GoodID:   "1",
				Quantity: 1,
			},
		},
		Status:       "Created",
		Name:         "name",
		FullName:     "test test",
		Address:      "via roma 1",
		Reservations: []string{},
		CreationTime: time.Now().UnixMilli(),
	}

	order, err := a.SendOrderUpdate(ctx, cmd)
	require.NoError(t, err)
	require.Equal(t, cmd.ID, order.ID)

	info, err := s.Info(ctx)
	require.NoError(t, err)

	require.Equal(t, info.State.Msgs, uint64(1))
}

func TestNatsStreamAdapterSendTransferUpdate(t *testing.T) {
	ctx := t.Context()

	ns, _ := broker.NewInProcessNATSServer(t)
	js, err := jetstream.New(ns)
	require.NoError(t, err)

	s, err := js.CreateStream(ctx, stream.TransferUpdateStreamConfig)
	require.NoError(t, err)

	broker, err := broker.NewNatsMessageBroker(ns)
	require.NoError(t, err)
	a := NewNatsStreamAdapter(broker)

	cmd := port.SendTransferUpdateCmd{
		ID: "1",
		Goods: []port.SendTransferUpdateGood{
			{
				GoodID:   "1",
				Quantity: 1,
			},
		},
		Status:        "Created",
		SenderID:      "1",
		ReceiverID:    "2",
		ReservationId: "",
		CreationTime:  time.Now().UnixMilli(),
	}

	tranfer, err := a.SendTransferUpdate(ctx, cmd)
	require.NoError(t, err)
	require.Equal(t, cmd.ID, tranfer.ID)

	info, err := s.Info(ctx)
	require.NoError(t, err)

	require.Equal(t, info.State.Msgs, uint64(1))
}

func TestNatsStreamAdapterSendContactOrder(t *testing.T) {
	ctx := t.Context()

	ns, _ := broker.NewInProcessNATSServer(t)
	js, err := jetstream.New(ns)
	require.NoError(t, err)

	s, err := js.CreateStream(ctx, internalStream.ContactWarehousesStreamConfig)
	require.NoError(t, err)

	broker, err := broker.NewNatsMessageBroker(ns)
	require.NoError(t, err)
	a := NewNatsStreamAdapter(broker)

	cmd := port.SendContactWarehouseCmd{
		Order: &model.Order{
			ID: "1",
			Goods: []model.GoodStock{
				{
					GoodID:   "1",
					Quantity: 1,
				},
			},
			Status:       "Created",
			Name:         "name",
			FullName:     "test test",
			Address:      "via roma 1",
			Reservations: []string{},
			CreationTime: time.Now().UnixMilli(),
		},
		Type:                  port.SendContactWarehouseTypeOrder,
		ConfirmedReservations: []port.ConfirmedReservation{},
		RetryInTime:           0,
		RetryUntil:            time.Now().Add(time.Hour).UnixMilli(),
	}

	err = a.SendContactWarehouses(ctx, cmd)
	require.NoError(t, err)

	info, err := s.Info(ctx)
	require.NoError(t, err)

	require.Equal(t, info.State.Msgs, uint64(1))
}

func TestNatsStreamAdapterSendContactTransfer(t *testing.T) {
	ctx := t.Context()

	ns, _ := broker.NewInProcessNATSServer(t)
	js, err := jetstream.New(ns)
	require.NoError(t, err)

	s, err := js.CreateStream(ctx, internalStream.ContactWarehousesStreamConfig)
	require.NoError(t, err)

	broker, err := broker.NewNatsMessageBroker(ns)
	require.NoError(t, err)
	a := NewNatsStreamAdapter(broker)

	cmd := port.SendContactWarehouseCmd{
		Transfer: &model.Transfer{
			ID: "1",
			Goods: []model.GoodStock{
				{
					GoodID:   "1",
					Quantity: 1,
				},
			},
			Status:            "Created",
			SenderID:          "1",
			ReceiverID:        "2",
			ReservationID:     "",
			LinkedStockUpdate: 0,
			CreationTime:      time.Now().UnixMilli(),
		},
		Type:                  port.SendContactWarehouseTypeTransfer,
		ConfirmedReservations: []port.ConfirmedReservation{},
		RetryInTime:           0,
		RetryUntil:            time.Now().Add(time.Hour).UnixMilli(),
	}

	err = a.SendContactWarehouses(ctx, cmd)
	require.NoError(t, err)

	info, err := s.Info(ctx)
	require.NoError(t, err)

	require.Equal(t, info.State.Msgs, uint64(1))
}

func TestNatsStreamAdapterRequestReservation(t *testing.T) {
	ctx := t.Context()

	ns, _ := broker.NewInProcessNATSServer(t)

	sub, err := ns.Subscribe("warehouse.1.reservation.create", func(msg *nats.Msg) {
		err := msg.Respond([]byte(`{"message": {"reservation_id":"1"}}`))
		require.NoError(t, err)
	})
	require.NoError(t, err)
	defer func() {
		err := sub.Unsubscribe()
		require.NoError(t, err)
	}()

	broker, err := broker.NewNatsMessageBroker(ns)
	require.NoError(t, err)
	a := NewNatsStreamAdapter(broker)

	cmd := port.RequestReservationCmd{
		WarehouseId: "1",
		Goods: []port.ReservationGood{
			{
				GoodID:   "1",
				Quantity: 1,
			},
		},
	}

	resp, err := a.RequestReservation(ctx, cmd)
	require.NoError(t, err)
	require.Equal(t, resp.Id, "1")
}
