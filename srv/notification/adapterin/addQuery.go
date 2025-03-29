package adapterin

import (
	"context"
	"encoding/json"
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/notification/portin"
	servicecmd "github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/nats-io/nats.go"
)

func NewAddQueryController(addQueryRuleUseCase portin.QueryRules) *AddQueryController {
	return &AddQueryController{addQueryRulePort: addQueryRuleUseCase}
}

type AddQueryController struct {
	addQueryRulePort portin.QueryRules
}

// Asserzione a compile-time che AddQueryController implementi Controller
var _ Controller = (*AddQueryController)(nil)

func (c *AddQueryController) Handle(ctx context.Context, msg *nats.Msg) error {
	request := stream.AddQueryRule{}
	err := json.Unmarshal(msg.Data, &request)
	if err != nil {
		_ = broker.RespondToMsg(msg, dto.InvalidJson())
		return err
	}

	cmd := servicecmd.QueryRule{GoodId: request.GoodID, Operator: request.Operator, Threshold: request.Threshold}
	id, err := c.addQueryRulePort.AddQueryRule(cmd)
	if err != nil {
		return err
	}

	return broker.RespondToMsg(msg, id.String())
}

func (c *AddQueryController) Subject() broker.Subject {
	return "notification.queries.add"
}
