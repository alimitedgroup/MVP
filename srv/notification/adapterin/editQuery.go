package adapterin

import (
	"context"
	"encoding/json"

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

var (
	EditQueryCounter metric.Int64Counter
)

func NewEditQueryController(addQueryRuleUseCase portin.QueryRules, mp MetricParams) *EditQueryController {
	observability.CounterSetup(&mp.Meter, mp.Logger, &TotalRequestCounter, &MetricMap, "num_notification_total_request")
	observability.CounterSetup(&mp.Meter, mp.Logger, &EditQueryCounter, &MetricMap, "num_notification_edit_query_request")
	Logger = mp.Logger
	return &EditQueryController{rulesPort: addQueryRuleUseCase}
}

type EditQueryController struct {
	rulesPort portin.QueryRules
}

// Asserzione a compile-time che EditQueryController implementi Controller
var _ Controller = (*EditQueryController)(nil)

func (c *EditQueryController) Handle(_ context.Context, msg *nats.Msg) error {

	Logger.Info("Received new edit query request")
	verdict := "success"

	defer func() {
		ctx := context.Background()
		Logger.Info("Edit query request terminated")
		TotalRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		EditQueryCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	var request dto.RuleEdit
	err := json.Unmarshal(msg.Data, &request)
	if err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		_ = broker.RespondToMsg(msg, dto.InvalidJson())
		return nil
	}

	cmd := types.EditRule{GoodId: request.GoodId, Operator: request.Operator, Threshold: request.Threshold}
	err = c.rulesPort.EditQueryRule(request.RuleId, cmd)
	if err != nil {
		verdict = "cannot handle request"
		Logger.Debug("Cannot handle request", zap.Error(err))
		_ = broker.RespondToMsg(msg, dto.InternalError())
	} else {
		_ = broker.RespondToMsg(msg, "OK")
	}

	return nil
}

func (c *EditQueryController) Subject() broker.Subject {
	return "notification.queries.add"
}
