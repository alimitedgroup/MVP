package adapterin

import (
	"context"
	"encoding/json"
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/notification/portin"
	"github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/nats-io/nats.go"
)

func NewEditQueryController(addQueryRuleUseCase portin.QueryRules) *EditQueryController {
	return &EditQueryController{rulesPort: addQueryRuleUseCase}
}

type EditQueryController struct {
	rulesPort portin.QueryRules
}

// Asserzione a compile-time che EditQueryController implementi Controller
var _ Controller = (*EditQueryController)(nil)

func (c *EditQueryController) Handle(_ context.Context, msg *nats.Msg) error {
	var request dto.RuleEdit
	err := json.Unmarshal(msg.Data, &request)
	if err != nil {
		_ = broker.RespondToMsg(msg, dto.InvalidJson())
		return nil
	}

	cmd := types.EditRule{GoodId: request.GoodId, Operator: request.Operator, Threshold: request.Threshold}
	err = c.rulesPort.EditQueryRule(request.RuleId, cmd)
	if err != nil {
		_ = broker.RespondToMsg(msg, dto.InternalError())
	} else {
		_ = broker.RespondToMsg(msg, "OK")
	}

	return nil
}

func (c *EditQueryController) Subject() broker.Subject {
	return "notification.queries.add"
}
