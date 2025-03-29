package adapterin

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/notification/portin"
	"github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

func NewGetQueryController(rulesPort portin.QueryRules) *GetQueryController {
	return &GetQueryController{rulesPort: rulesPort}
}

type GetQueryController struct {
	rulesPort portin.QueryRules
}

// Asserzione a compile-time che GetQueryController implementi Controller
var _ Controller = (*GetQueryController)(nil)

func (c *GetQueryController) Handle(_ context.Context, msg *nats.Msg) error {
	var request string
	err := json.Unmarshal(msg.Data, &request)
	if err != nil {
		_ = broker.RespondToMsg(msg, dto.InvalidJson())
		return nil
	}

	id, err := uuid.Parse(request)
	if err != nil {
		_ = broker.RespondToMsg(msg, dto.InvalidJson())
		return nil
	}

	rule, err := c.rulesPort.GetQueryRule(id)
	if errors.Is(err, types.ErrRuleNotExists) {
		_ = broker.RespondToMsg(msg, dto.RuleNotFound())
	} else if err != nil {
		_ = broker.RespondToMsg(msg, dto.InternalError())
	} else {
		_ = broker.RespondToMsg(msg, rule)
	}

	return nil
}

func (c *GetQueryController) Subject() broker.Subject {
	return "notification.queries.get"
}
