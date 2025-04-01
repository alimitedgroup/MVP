package adapterin

import (
	"context"
	"encoding/json"
	"go.uber.org/fx"
	"time"

	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/notification/portin"
	servicecmd "github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

var (
	StockUpdateCounter metric.Int64Counter
)

type StockReceiverParams struct {
	fx.In

	StockPort portin.StockUpdates
	Logger    *zap.Logger
	Meter     metric.Meter
}

func NewStockReceiver(mp StockReceiverParams) *StockUpdateReceiver {
	observability.CounterSetup(&mp.Meter, mp.Logger, &TotalRequestCounter, &MetricMap, "num_notification_total_request")
	observability.CounterSetup(&mp.Meter, mp.Logger, &StockUpdateCounter, &MetricMap, "num_notification_stock_update_query_request")
	return &StockUpdateReceiver{
		stockPort: mp.StockPort,
		Logger:    mp.Logger,
	}
}

type StockUpdateReceiver struct {
	stockPort portin.StockUpdates
	*zap.Logger
}

var _ JsController = (*StockUpdateReceiver)(nil)

func (s StockUpdateReceiver) Handle(_ context.Context, msg jetstream.Msg) error {
	s.Info("Received new stock update query request")
	verdict := "success"

	defer func() {
		ctx := context.Background()
		s.Info("Stock update query request terminated")
		TotalRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		StockUpdateCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	request := &stream.StockUpdate{}

	err := json.Unmarshal(msg.Data(), request)
	if err != nil {
		verdict = "bad request"
		s.Debug("Bad request", zap.Error(err))
		return err
	}

	goods := make([]servicecmd.StockGood, len(request.Goods))
	for i, g := range request.Goods {
		goods[i] = servicecmd.StockGood{
			ID:       g.GoodID,
			Quantity: int(g.Quantity),
			Delta:    int(g.Delta),
		}
	}

	cmd := servicecmd.NewAddStockUpdateCmd(request.WarehouseID, string(request.Type), request.OrderID, request.TransferID, goods, time.Now().Unix())
	return s.stockPort.RecordStockUpdate(cmd)
}

func (s StockUpdateReceiver) Stream() jetstream.StreamConfig {
	return stream.StockUpdateStreamConfig
}
