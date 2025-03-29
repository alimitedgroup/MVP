package listener

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/order/adapter/controller"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	StockUpdateCounter metric.Int64Counter
	Logger             *zap.Logger
)

type MetricParams struct {
	fx.In
	Logger *zap.Logger
	Meter  metric.Meter
}

type StockUpdateListener struct {
	applyStockUpdateUseCase port.IApplyStockUpdateUseCase
}

func NewStockUpdateListener(applyStockUpdateUseCase port.IApplyStockUpdateUseCase, mp MetricParams) *StockUpdateListener {
	observability.CounterSetup(&mp.Meter, mp.Logger, &StockUpdateCounter, &controller.MetricMap, "num_stock_updates_requests")
	observability.CounterSetup(&mp.Meter, mp.Logger, &controller.TotalRequestCounter, &controller.MetricMap, "num_order_transfer_requests")
	Logger = mp.Logger
	return &StockUpdateListener{applyStockUpdateUseCase}
}

func (l *StockUpdateListener) ListenStockUpdate(ctx context.Context, msg jetstream.Msg) error {

	Logger.Info("Received stock update request")
	verdict := "success"

	defer func() {
		Logger.Info("Stock update request terminated")
		StockUpdateCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		controller.TotalRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	var event stream.StockUpdate

	err := json.Unmarshal(msg.Data(), &event)
	if err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		return err
	}

	cmd := stockUpdateEventToApplyStockUpdateCmd(event)
	err = l.applyStockUpdateUseCase.ApplyStockUpdate(ctx, cmd)
	if err != nil {
		verdict = "cannot apply stock update"
		Logger.Debug("Cannot apply stock update", zap.Error(err))
		return err
	}

	return nil
}

func stockUpdateEventToApplyStockUpdateCmd(event stream.StockUpdate) port.StockUpdateCmd {
	goods := make([]port.StockUpdateGood, 0, len(event.Goods))

	for _, good := range event.Goods {
		goods = append(goods, port.StockUpdateGood{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
			Delta:    good.Delta,
		})
	}

	return port.StockUpdateCmd{
		ID:            event.ID,
		WarehouseID:   event.WarehouseID,
		Type:          port.StockUpdateType(event.Type),
		OrderID:       event.OrderID,
		TransferID:    event.TransferID,
		ReservationID: event.ReservationID,
		Timestamp:     event.Timestamp,
		Goods:         goods,
	}
}
