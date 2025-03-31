package listener

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/controller"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

var (
	ReservationEventCounter metric.Int64Counter
)

type ReservationEventListener struct {
	applyReservationEventUseCase port.IApplyReservationUseCase
}

func NewReservationEventListener(applyReservationEventUseCase port.IApplyReservationUseCase, mp MetricParams) *ReservationEventListener {
	observability.CounterSetup(&mp.Meter, mp.Logger, &ReservationEventCounter, &controller.MetricMap, "num_reservation_event_requests")
	observability.CounterSetup(&mp.Meter, mp.Logger, &controller.TotalRequestsCounter, &controller.MetricMap, "num_warehouse_requests")
	Logger = mp.Logger
	return &ReservationEventListener{applyReservationEventUseCase}
}

func (l *ReservationEventListener) ListenReservationEvent(ctx context.Context, msg jetstream.Msg) error {

	Logger.Info("Received reservation event request")
	verdict := "success"

	defer func() {
		Logger.Info("Reservation event request terminated")
		ReservationEventCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		controller.TotalRequestsCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	var event stream.ReservationEvent
	err := json.Unmarshal(msg.Data(), &event)
	if err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		return err
	}
	Logger.Debug("Reservation event request", zap.Any("event", event))

	goods := make([]port.ReservationGood, 0, len(event.Goods))
	for _, good := range event.Goods {
		goods = append(goods, port.ReservationGood(good))
	}
	cmd := port.ApplyReservationEventCmd{
		ID:    event.ID,
		Goods: goods,
	}
	err = l.applyReservationEventUseCase.ApplyReservationEvent(cmd)
	if err != nil {
		verdict = "cannot apply resevation event"
		Logger.Debug("Cannot apply reservation event", zap.Error(err))
		return err
	}

	return nil
}
