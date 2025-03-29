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
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	GoodUpdateRequest metric.Int64Counter
	Logger            *zap.Logger
)

type MetricParams struct {
	fx.In
	Logger *zap.Logger
	Meter  metric.Meter
}

type CatalogListener struct {
	applyCatalogUpdateUseCase port.IApplyCatalogUpdateUseCase
}

func NewCatalogListener(applyCatalogUpdateUseCase port.IApplyCatalogUpdateUseCase, mp MetricParams) *CatalogListener {
	observability.CounterSetup(&mp.Meter, mp.Logger, &GoodUpdateRequest, &controller.MetricMap, "num_update_good_requests")
	observability.CounterSetup(&mp.Meter, mp.Logger, &controller.TotalRequestsCounter, &controller.MetricMap, "num_warehouse_requests")
	Logger = mp.Logger
	return &CatalogListener{applyCatalogUpdateUseCase}
}

func (l *CatalogListener) ListenGoodUpdate(ctx context.Context, msg jetstream.Msg) error {

	Logger.Info("Received good update request")
	verdict := "success"

	defer func() {
		Logger.Info("Good update request terminated")
		GoodUpdateRequest.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		controller.TotalRequestsCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	var event stream.GoodUpdateData
	err := json.Unmarshal(msg.Data(), &event)
	if err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		return err
	}

	cmd := port.CatalogUpdateCmd{
		GoodID:      event.GoodID,
		Name:        event.GoodNewName,
		Description: event.GoodNewDescription,
	}

	l.applyCatalogUpdateUseCase.ApplyCatalogUpdate(cmd)

	return nil
}
