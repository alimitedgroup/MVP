package listener

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/stream"
	internalStream "github.com/alimitedgroup/MVP/srv/order/adapter/stream"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/nats-io/nats.go/jetstream"
)

type OrderListener struct {
	applyOrderUpdateUseCase port.IApplyOrderUpdateUseCase
	contactWarehouseUseCase port.IContactWarehousesUseCase
}

func NewOrderListener(applyOrderUpdateUseCase port.IApplyOrderUpdateUseCase, contactWarehouseUseCase port.IContactWarehousesUseCase) *OrderListener {
	return &OrderListener{applyOrderUpdateUseCase, contactWarehouseUseCase}
}

func (l *OrderListener) ListenOrderUpdate(ctx context.Context, msg jetstream.Msg) error {
	var event stream.OrderUpdate
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		return err
	}

	cmd := orderUpdateEventToApplyOrderUpdateCmd(event)
	if err := l.applyOrderUpdateUseCase.ApplyOrderUpdate(ctx, cmd); err != nil {
		return err
	}

	return nil
}

func (l *OrderListener) ListenContactWarehouses(ctx context.Context, msg jetstream.Msg) error {
	var event internalStream.ContactWarehouses
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		return err
	}

	confirmed := make([]port.ConfirmedReservation, 0, len(event.ConfirmedReservations))
	for _, reservation := range event.ConfirmedReservations {
		confirmed = append(confirmed, port.ConfirmedReservation{
			WarehouseId:   reservation.WarehouseId,
			ReservationID: reservation.ReservationID,
			Goods:         reservation.Goods,
		})
	}

	goods := make([]port.ContactWarehousesGood, 0, len(event.Order.Goods))

	for _, good := range event.Order.Goods {
		goods = append(goods, port.ContactWarehousesGood{
			GoodId:   good.GoodId,
			Quantity: good.Quantity,
		})
	}

	cmd := port.ContactWarehousesCmd{
		Order: port.ContactWarehousesOrder{
			ID:           event.Order.ID,
			Goods:        goods,
			Status:       event.Order.Status,
			Name:         event.Order.Name,
			FullName:     event.Order.FullName,
			Address:      event.Order.Address,
			UpdateTime:   event.Order.UpdateTime,
			CreationTime: event.Order.CreationTime,
			Reservations: event.Order.Reservations,
		},
		LastContact:           event.LastContact,
		ConfirmedReservations: confirmed,
		ExcludeWarehouses:     event.ExcludeWarehouses,
	}
	err := l.contactWarehouseUseCase.ContactWarehouses(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}

func orderUpdateEventToApplyOrderUpdateCmd(event stream.OrderUpdate) port.OrderUpdateCmd {
	goods := make([]port.OrderUpdateGood, 0, len(event.Goods))
	for _, good := range event.Goods {
		goods = append(goods, port.OrderUpdateGood{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}

	cmd := port.OrderUpdateCmd{
		ID:           event.ID,
		Goods:        goods,
		Status:       event.Status,
		Name:         event.Name,
		FullName:     event.Email,
		Address:      event.Address,
		Reservations: event.Reservations,
		UpdateTime:   event.UpdateTime,
		CreationTime: event.CreationTime,
	}

	return cmd
}
