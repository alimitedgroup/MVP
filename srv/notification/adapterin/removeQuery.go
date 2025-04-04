package adapterin

import (
	"context"
	"errors"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/srv/notification/portin"
	"github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

var (
	RemoveQueryCounter metric.Int64Counter
)

func NewRemoveQueryController(p QueryControllersParams) *RemoveQueryController {
	observability.CounterSetup(&p.Meter, p.Logger, &TotalRequestCounter, &MetricMap, "num_notification_total_request")
	observability.CounterSetup(&p.Meter, p.Logger, &RemoveQueryCounter, &MetricMap, "num_notification_remove_query_request")
	return &RemoveQueryController{rulesPort: p.RulesPort, Logger: p.Logger}
}

type RemoveQueryController struct {
	rulesPort portin.QueryRules
	*zap.Logger
}

// Asserzione a compile-time che AddQueryController implementi Controller
var _ Controller = (*RemoveQueryController)(nil)

func (c *RemoveQueryController) Handle(_ context.Context, msg *nats.Msg) error {
	c.Info("Received new remove query request")
	verdict := "success"

	defer func() {
		ctx := context.Background()
		c.Info("Remove query request terminated")
		TotalRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		ListQueryCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	request := string(msg.Data)

	id, err := uuid.Parse(request)
	if err != nil {
		verdict = "bad request"
		c.Debug("Bad request", zap.Error(err))
		_ = broker.RespondToMsg(msg, dto.InvalidJson())
		return nil
	}

	err = c.rulesPort.RemoveQueryRule(id)
	if errors.Is(err, types.ErrRuleNotExists) {
		verdict = "cannot handle request"
		c.Debug("Cannot handle request", zap.Error(err))
		_ = broker.RespondToMsg(msg, dto.RuleNotFound())
		return nil
	}
	if err != nil {
		_ = broker.RespondToMsg(msg, dto.InternalError())
		return nil
	}

	_ = msg.Respond([]byte("OK"))
	return nil
}

func (c *RemoveQueryController) Subject() broker.Subject {
	return "notification.queries.remove"
}
