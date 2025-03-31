package listener

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/controller"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

var (
	OrderUpdateRequestCounter    metric.Int64Counter
	TransferUpdateRequestCounter metric.Int64Counter
)

type OrderUpdateListener struct {
	confirmOrderUseCase    port.IConfirmOrderUseCase
	confirmTransferUseCase port.IConfirmTransferUseCase
}

func NewOrderUpdateListener(confirmOrderUseCase port.IConfirmOrderUseCase, confirmTransferUseCase port.IConfirmTransferUseCase, mp MetricParams) *OrderUpdateListener {
	observability.CounterSetup(&mp.Meter, mp.Logger, &OrderUpdateRequestCounter, &controller.MetricMap, "num_update_order_requests")
	observability.CounterSetup(&mp.Meter, mp.Logger, &TransferUpdateRequestCounter, &controller.MetricMap, "num_update_transfer_requests")
	observability.CounterSetup(&mp.Meter, mp.Logger, &controller.TotalRequestsCounter, &controller.MetricMap, "num_warehouse_requests")
	Logger = mp.Logger
	return &OrderUpdateListener{confirmOrderUseCase, confirmTransferUseCase}
}

func (l *OrderUpdateListener) ListenOrderUpdate(ctx context.Context, msg jetstream.Msg) error {
	Logger.Info("Received order update request")
	verdict := "success"

	defer func() {
		Logger.Info("Order update request terminated")
		OrderUpdateRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		controller.TotalRequestsCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	var event stream.OrderUpdate
	err := json.Unmarshal(msg.Data(), &event)
	if err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		return err
	}
	Logger.Debug("Order update event", zap.Any("event", event))

	goods := make([]port.OrderUpdateGood, 0, len(event.Goods))
	for _, good := range event.Goods {
		goods = append(goods, port.OrderUpdateGood(good))
	}
	cmd := port.ConfirmOrderCmd{
		OrderID:      event.ID,
		Status:       event.Status,
		Reservations: event.Reservations,
		Goods:        goods,
	}
	err = l.confirmOrderUseCase.ConfirmOrder(ctx, cmd)
	if err != nil {
		verdict = "cannot confirm order"
		Logger.Debug("Cannot confirm order", zap.Error(err))
		return err
	}

	return nil
}

func (l *OrderUpdateListener) ListenTransferUpdate(ctx context.Context, msg jetstream.Msg) error {

	Logger.Info("Received transfer update request")
	verdict := "success"

	defer func() {
		Logger.Info("Transfer update request terminated")
		TransferUpdateRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		controller.TotalRequestsCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	var event stream.TransferUpdate
	err := json.Unmarshal(msg.Data(), &event)
	if err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		return err
	}

	goods := make([]port.TransferUpdateGood, 0, len(event.Goods))
	for _, good := range event.Goods {
		goods = append(goods, port.TransferUpdateGood(good))
	}
	cmd := port.ConfirmTransferCmd{
		TransferID:    event.ID,
		Status:        event.Status,
		SenderID:      event.SenderID,
		ReceiverID:    event.ReceiverID,
		ReservationID: event.ReservationID,
		Goods:         goods,
	}
	err = l.confirmTransferUseCase.ConfirmTransfer(ctx, cmd)
	if err != nil {
		verdict = "cannot confirm transfer"
		Logger.Debug("Cannot confirm transfer", zap.Error(err))
		return err
	}

	return nil
}
