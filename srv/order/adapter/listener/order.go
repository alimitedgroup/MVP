package listener

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	internalStream "github.com/alimitedgroup/MVP/srv/order/adapter/stream"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/nats-io/nats.go/jetstream"
)

type OrderListener struct {
	applyOrderUpdateUseCase    port.IApplyOrderUpdateUseCase
	applyTransferUpdateUseCase port.IApplyTransferUpdateUseCase
	contactWarehouseUseCase    port.IContactWarehousesUseCase
}

func NewOrderListener(applyOrderUpdateUseCase port.IApplyOrderUpdateUseCase, contactWarehouseUseCase port.IContactWarehousesUseCase, applyTransferUpdateUseCase port.IApplyTransferUpdateUseCase) *OrderListener {
	return &OrderListener{applyOrderUpdateUseCase, applyTransferUpdateUseCase, contactWarehouseUseCase}
}

func (l *OrderListener) ListenOrderUpdate(ctx context.Context, msg jetstream.Msg) error {
	var event stream.OrderUpdate
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		return err
	}

	cmd := orderUpdateEventToApplyOrderUpdateCmd(event)
	l.applyOrderUpdateUseCase.ApplyOrderUpdate(ctx, cmd)

	return nil
}

func (l *OrderListener) ListenTransferUpdate(ctx context.Context, msg jetstream.Msg) error {
	var event stream.TransferUpdate
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		return err
	}

	cmd := transferUpdateEventToApplyTransferUpdateCmd(event)
	l.applyTransferUpdateUseCase.ApplyTransferUpdate(ctx, cmd)

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

	var order *port.ContactWarehousesOrder
	var transfer *port.ContactWarehousesTransfer
	if event.Type == internalStream.ContactWarehousesTypeOrder {
		goods := make([]port.ContactWarehousesGood, 0, len(event.Order.Goods))

		for _, good := range event.Order.Goods {
			goods = append(goods, port.ContactWarehousesGood{
				GoodID:   good.GoodID,
				Quantity: good.Quantity,
			})
		}
		order = &port.ContactWarehousesOrder{
			ID:           event.Order.ID,
			Goods:        goods,
			Status:       event.Order.Status,
			Name:         event.Order.Name,
			FullName:     event.Order.FullName,
			Address:      event.Order.Address,
			UpdateTime:   event.Order.UpdateTime,
			CreationTime: event.Order.CreationTime,
			Reservations: event.Order.Reservations,
		}
	} else if event.Type == internalStream.ContactWarehousesTypeTransfer {
		goods := make([]port.ContactWarehousesGood, 0, len(event.Transfer.Goods))

		for _, good := range event.Transfer.Goods {
			goods = append(goods, port.ContactWarehousesGood{
				GoodID:   good.GoodID,
				Quantity: good.Quantity,
			})
		}

		transfer = &port.ContactWarehousesTransfer{
			ID:            event.Transfer.ID,
			Goods:         goods,
			SenderID:      event.Transfer.SenderId,
			ReceiverID:    event.Transfer.ReceiverId,
			Status:        event.Transfer.Status,
			UpdateTime:    event.Transfer.UpdateTime,
			CreationTime:  event.Transfer.CreationTime,
			ReservationId: event.Transfer.ReservationId,
		}
	}

	cmd := port.ContactWarehousesCmd{
		Type:                  port.ContactWarehousesType(event.Type),
		Order:                 order,
		Transfer:              transfer,
		ConfirmedReservations: confirmed,
		ExcludeWarehouses:     event.ExcludeWarehouses,
		RetryInTime:           event.RetryInTime,
		RetryUntil:            event.RetryUntil,
	}

	retry, err := l.contactWarehouseUseCase.ContactWarehouses(ctx, cmd)
	if err != nil {
		return err
	}

	if retry.IsRetry {
		if err := msg.NakWithDelay(retry.RetryAfter); err != nil {
			return err
		}
		return broker.ErrMsgNotAcked
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
		FullName:     event.FullName,
		Address:      event.Address,
		Reservations: event.Reservations,
		UpdateTime:   event.UpdateTime,
		CreationTime: event.CreationTime,
	}

	return cmd
}

func transferUpdateEventToApplyTransferUpdateCmd(event stream.TransferUpdate) port.TransferUpdateCmd {
	goods := make([]port.TransferUpdateGood, 0, len(event.Goods))
	for _, good := range event.Goods {
		goods = append(goods, port.TransferUpdateGood{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}

	cmd := port.TransferUpdateCmd{
		ID:            event.ID,
		Goods:         goods,
		SenderId:      event.SenderID,
		ReceiverId:    event.ReceiverID,
		Status:        event.Status,
		ReservationId: event.ReservationId,
		UpdateTime:    event.UpdateTime,
		CreationTime:  event.CreationTime,
	}

	return cmd
}
