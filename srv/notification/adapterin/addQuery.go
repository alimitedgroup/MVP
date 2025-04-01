package adapterin

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/srv/notification/portin"
	"github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	MetricMap           sync.Map
	TotalRequestCounter metric.Int64Counter
	AddQueryCounter     metric.Int64Counter
)

type AddQueryParams struct {
	fx.In
	Logger    *zap.Logger
	Meter     metric.Meter
	rulesPort portin.QueryRules
}

func NewAddQueryController(p AddQueryParams) *AddQueryController {
	observability.CounterSetup(&p.Meter, p.Logger, &TotalRequestCounter, &MetricMap, "num_notification_total_request")
	observability.CounterSetup(&p.Meter, p.Logger, &AddQueryCounter, &MetricMap, "num_notification_add_query_request")
	return &AddQueryController{rulesPort: p.rulesPort, Logger: p.Logger}
}

type AddQueryController struct {
	rulesPort portin.QueryRules
	*zap.Logger
}

// Asserzione a compile-time che AddQueryController implementi Controller
var _ Controller = (*AddQueryController)(nil)

func (c *AddQueryController) Handle(_ context.Context, msg *nats.Msg) error {
	c.Info("Received new add query request")
	verdict := "success"

	defer func() {
		ctx := context.Background()
		c.Info("Add query request terminated")
		TotalRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		AddQueryCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	var request dto.Rule
	err := json.Unmarshal(msg.Data, &request)
	if err != nil {
		verdict = "bad request"
		c.Debug("Bad request", zap.Error(err))
		_ = broker.RespondToMsg(msg, dto.InvalidJson())
		return nil
	}

	cmd := types.QueryRule{GoodId: request.GoodId, Operator: request.Operator, Threshold: request.Threshold}
	id, err := c.rulesPort.AddQueryRule(cmd)
	if err != nil {
		verdict = "cannot handle request"
		c.Debug("Cannot handle request", zap.Error(err))
		_ = broker.RespondToMsg(msg, dto.InternalError())
	} else {
		_ = msg.Respond([]byte(id.String()))
	}

	return nil
}

func (c *AddQueryController) Subject() broker.Subject {
	return "notification.queries.add"
}
