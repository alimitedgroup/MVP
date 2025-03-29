package listener

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/order/adapter/controller"
	internalStream "github.com/alimitedgroup/MVP/srv/order/adapter/stream"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	OrderUpdateRequestCounter      metric.Int64Counter
	TransferUpdateRequestCounter   metric.Int64Counter
	ContactWarehouseRequestCounter metric.Int64Counter
)

type OrderListener struct {
	applyOrderUpdateUseCase    port.IApplyOrderUpdateUseCase
	applyTransferUpdateUseCase port.IApplyTransferUpdateUseCase
	contactWarehouseUseCase    port.IContactWarehousesUseCase
}

type OrderListenerParams struct {
	fx.In

	ApplyOrderUpdateUseCase    port.IApplyOrderUpdateUseCase
	ApplyTransferUpdateUseCase port.IApplyTransferUpdateUseCase
	ContactWarehouseUseCase    port.IContactWarehousesUseCase
	Logger                     *zap.Logger
	Meter                      metric.Meter
}

func NewOrderListener(p OrderListenerParams) *OrderListener {
	observability.CounterSetup(&p.Meter, p.Logger, &OrderUpdateRequestCounter, &controller.MetricMap, "num_order_updates_requests")
	observability.CounterSetup(&p.Meter, p.Logger, &TransferUpdateRequestCounter, &controller.MetricMap, "num_transfer_updates_requests")
	observability.CounterSetup(&p.Meter, p.Logger, &ContactWarehouseRequestCounter, &controller.MetricMap, "num_contact_warehouse_requests")
	observability.CounterSetup(&p.Meter, p.Logger, &controller.TotalRequestCounter, &controller.MetricMap, "num_order_transfer_requests")
	Logger = p.Logger
	return &OrderListener{p.ApplyOrderUpdateUseCase, p.ApplyTransferUpdateUseCase, p.ContactWarehouseUseCase}
}

func (l *OrderListener) ListenOrderUpdate(ctx context.Context, msg jetstream.Msg) error {

	Logger.Info("Received good update request")
	verdict := "success"

	defer func() {
		Logger.Info("Good update request terminated")
		OrderUpdateRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		controller.TotalRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	var event stream.OrderUpdate
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		return err
	}

	cmd := orderUpdateEventToApplyOrderUpdateCmd(event)
	l.applyOrderUpdateUseCase.ApplyOrderUpdate(ctx, cmd)

	return nil
}

func (l *OrderListener) ListenTransferUpdate(ctx context.Context, msg jetstream.Msg) error {

	Logger.Info("Received transfer update request")
	verdict := "success"

	defer func() {
		Logger.Info("Transfer update request terminated")
		OrderUpdateRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		controller.TotalRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	var event stream.TransferUpdate
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		return err
	}

	cmd := transferUpdateEventToApplyTransferUpdateCmd(event)
	l.applyTransferUpdateUseCase.ApplyTransferUpdate(ctx, cmd)

	return nil
}

func (l *OrderListener) ListenContactWarehouses(ctx context.Context, msg jetstream.Msg) error {

	Logger.Info("Received contact warehouse request")
	verdict := "success"

	defer func() {
		Logger.Info("Contact warehouse request terminated")
		ContactWarehouseRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		controller.TotalRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	var event internalStream.ContactWarehouses
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		return err
	}

	confirmed := make([]port.ConfirmedReservation, 0, len(event.ConfirmedReservations))
	for _, reservation := range event.ConfirmedReservations {
		confirmed = append(confirmed, port.ConfirmedReservation{
			WarehouseID:   reservation.WarehouseID,
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
			SenderID:      event.Transfer.SenderID,
			ReceiverID:    event.Transfer.ReceiverID,
			Status:        event.Transfer.Status,
			UpdateTime:    event.Transfer.UpdateTime,
			CreationTime:  event.Transfer.CreationTime,
			ReservationID: event.Transfer.ReservationID,
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
		verdict = "cannot contact warehouse"
		Logger.Debug("Cannot contact warehouse", zap.Error(err))
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
		SenderID:      event.SenderID,
		ReceiverID:    event.ReceiverID,
		Status:        event.Status,
		ReservationID: event.ReservationID,
		UpdateTime:    event.UpdateTime,
		CreationTime:  event.CreationTime,
	}

	return cmd
}
