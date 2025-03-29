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

func NewAddQueryController(addQueryRuleUseCase portin.QueryRules) *AddQueryController {
	return &AddQueryController{rulesPort: addQueryRuleUseCase}
}

type AddQueryController struct {
	rulesPort portin.QueryRules
}

// Asserzione a compile-time che AddQueryController implementi Controller
var _ Controller = (*AddQueryController)(nil)

func (c *AddQueryController) Handle(_ context.Context, msg *nats.Msg) error {
	var request dto.Rule
	err := json.Unmarshal(msg.Data, &request)
	if err != nil {
		_ = broker.RespondToMsg(msg, dto.InvalidJson())
		return nil
	}

	cmd := types.QueryRule{GoodId: request.GoodId, Operator: request.Operator, Threshold: request.Threshold}
	id, err := c.rulesPort.AddQueryRule(cmd)
	if err != nil {
		_ = broker.RespondToMsg(msg, dto.InternalError())
	} else {
		_ = broker.RespondToMsg(msg, id.String())
	}

	return nil
}

func (c *AddQueryController) Subject() broker.Subject {
	return "notification.queries.add"
}
