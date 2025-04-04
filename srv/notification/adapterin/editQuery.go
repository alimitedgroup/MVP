package adapterin

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/srv/notification/portin"
	"github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

var EditQueryCounter metric.Int64Counter

type QueryControllersParams struct {
	fx.In

	RulesPort portin.QueryRules
	Logger    *zap.Logger
	Meter     metric.Meter
}

func NewEditQueryController(p QueryControllersParams) *EditQueryController {
	observability.CounterSetup(&p.Meter, p.Logger, &TotalRequestCounter, &MetricMap, "num_notification_total_request")
	observability.CounterSetup(&p.Meter, p.Logger, &EditQueryCounter, &MetricMap, "num_notification_edit_query_request")
	return &EditQueryController{rulesPort: p.RulesPort, Logger: p.Logger}
}

type EditQueryController struct {
	rulesPort portin.QueryRules
	*zap.Logger
}

// Asserzione a compile-time che EditQueryController implementi Controller
var _ Controller = (*EditQueryController)(nil)

func (c *EditQueryController) Handle(_ context.Context, msg *nats.Msg) error {
	c.Info("Received new edit query request")
	verdict := "success"

	defer func() {
		ctx := context.Background()
		c.Info("Edit query request terminated")
		TotalRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		EditQueryCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	var request dto.RuleEdit
	err := json.Unmarshal(msg.Data, &request)
	if err != nil {
		verdict = "bad request"
		c.Debug("Bad request", zap.Error(err))
		_ = broker.RespondToMsg(msg, dto.InvalidJson())
		return nil
	}

	cmd := types.EditRule{GoodId: request.GoodId, Operator: request.Operator, Threshold: request.Threshold}
	err = c.rulesPort.EditQueryRule(request.RuleId, cmd)
	if errors.Is(err, jetstream.ErrKeyNotFound) {
		verdict = "not found"
		_ = broker.RespondToMsg(msg, dto.RuleNotFound())
	} else if err != nil {
		verdict = "cannot handle request"
		c.Debug("Cannot handle request", zap.Error(err))
		_ = broker.RespondToMsg(msg, dto.InternalError())
	} else {
		_ = msg.Respond([]byte("OK"))
	}

	return nil
}

func (c *EditQueryController) Subject() broker.Subject {
	return "notification.queries.edit"
}
