package adapterin

import (
	"context"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/srv/notification/portin"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

var (
	ListQueryCounter metric.Int64Counter
)

func NewListQueriesController(rulesPort portin.QueryRules, mp AddQueryParams) *ListQueriesController {
	observability.CounterSetup(&mp.Meter, mp.Logger, &TotalRequestCounter, &MetricMap, "num_notification_total_request")
	observability.CounterSetup(&mp.Meter, mp.Logger, &ListQueryCounter, &MetricMap, "num_notification_list_query_request")
	return &ListQueriesController{rulesPort: rulesPort, Logger: mp.Logger}
}

type ListQueriesController struct {
	rulesPort portin.QueryRules
	*zap.Logger
}

// Asserzione a compile-time che ListQueriesController implementi Controller
var _ Controller = (*ListQueriesController)(nil)

func (c *ListQueriesController) Handle(_ context.Context, msg *nats.Msg) error {
	c.Info("Received new list query request")
	verdict := "success"

	defer func() {
		ctx := context.Background()
		c.Info("List query request terminated")
		TotalRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		ListQueryCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	rules, err := c.rulesPort.ListQueryRules()
	if err != nil {
		verdict = "cannot handle request"
		c.Debug("Cannot handle request", zap.Error(err))
		_ = broker.RespondToMsg(msg, dto.InternalError())
	} else {
		_ = broker.RespondToMsg(msg, rules)
	}

	return nil
}

func (c *ListQueriesController) Subject() broker.Subject {
	return "notification.queries.list"
}
