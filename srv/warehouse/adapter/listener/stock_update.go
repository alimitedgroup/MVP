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
	StockUpdateCounter metric.Int64Counter
)

type StockUpdateListener struct {
	applyStockUpdateUseCase port.IApplyStockUpdateUseCase
}

func NewStockUpdateListener(applyStockUpdateUseCase port.IApplyStockUpdateUseCase, mp MetricParams) *StockUpdateListener {
	observability.CounterSetup(&mp.Meter, mp.Logger, &StockUpdateCounter, &controller.MetricMap, "num_update_stock_requests")
	observability.CounterSetup(&mp.Meter, mp.Logger, &controller.TotalRequestsCounter, &controller.MetricMap, "num_warehouse_requests")
	Logger = mp.Logger
	return &StockUpdateListener{applyStockUpdateUseCase}
}

func (l *StockUpdateListener) ListenStockUpdate(ctx context.Context, msg jetstream.Msg) error {

	Logger.Info("Received stock update request")
	verdict := "success"

	defer func() {
		Logger.Info("Stock update request terminated")
		StockUpdateCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		controller.TotalRequestsCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	var event stream.StockUpdate
	err := json.Unmarshal(msg.Data(), &event)
	if err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		return err
	}
	Logger.Debug("Stock update request", zap.Any("event", event))

	cmd := stockUpdateEventToApplyStockUpdateCmd(event)
	l.applyStockUpdateUseCase.ApplyStockUpdate(cmd)

	return nil
}

func stockUpdateEventToApplyStockUpdateCmd(event stream.StockUpdate) port.StockUpdateCmd {
	goods := make([]port.StockUpdateGood, 0, len(event.Goods))
	for _, good := range event.Goods {
		goods = append(goods, port.StockUpdateGood(good))
	}
	return port.StockUpdateCmd{
		ID:            event.ID,
		Type:          port.StockUpdateType(event.Type),
		OrderID:       event.OrderID,
		TransferID:    event.TransferID,
		ReservationID: event.ReservationID,
		Timestamp:     event.Timestamp,
		Goods:         goods,
	}
}
