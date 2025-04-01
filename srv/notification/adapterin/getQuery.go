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
	GetQueryCounter metric.Int64Counter
)

func NewGetQueryController(rulesPort portin.QueryRules, mp AddQueryParams) *GetQueryController {
	observability.CounterSetup(&mp.Meter, mp.Logger, &TotalRequestCounter, &MetricMap, "num_notification_total_request")
	observability.CounterSetup(&mp.Meter, mp.Logger, &GetQueryCounter, &MetricMap, "num_notification_get_query_request")
	return &GetQueryController{rulesPort: rulesPort, Logger: mp.Logger}
}

type GetQueryController struct {
	rulesPort portin.QueryRules
	*zap.Logger
}

// Asserzione a compile-time che GetQueryController implementi Controller
var _ Controller = (*GetQueryController)(nil)

func (c *GetQueryController) Handle(_ context.Context, msg *nats.Msg) error {

	c.Info("Received new get query request")
	verdict := "success"

	defer func() {
		ctx := context.Background()
		c.Info("Get query request terminated")
		TotalRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		GetQueryCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	request := string(msg.Data)
	id, err := uuid.Parse(request)
	if err != nil {
		verdict = "bad request"
		c.Debug("Bad request", zap.Error(err))
		_ = broker.RespondToMsg(msg, dto.InvalidJson())
		return nil
	}

	rule, err := c.rulesPort.GetQueryRule(id)
	if errors.Is(err, types.ErrRuleNotExists) {
		verdict = "cannot handle request"
		c.Debug("Cannot handle request", zap.Error(err))
		_ = broker.RespondToMsg(msg, dto.RuleNotFound())
	} else if err != nil {
		verdict = "cannot handle request"
		c.Debug("Cannot handle request", zap.Error(err))
		_ = broker.RespondToMsg(msg, dto.InternalError())
	} else {
		_ = broker.RespondToMsg(msg, rule)
	}

	return nil
}

func (c *GetQueryController) Subject() broker.Subject {
	return "notification.queries.get"
}
