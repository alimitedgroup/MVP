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

func NewRemoveQueryController(addQueryRuleUseCase portin.QueryRules) *RemoveQueryController {
	return &RemoveQueryController{rulesPort: addQueryRuleUseCase}
}

type RemoveQueryController struct {
	rulesPort portin.QueryRules
}

// Asserzione a compile-time che AddQueryController implementi Controller
var _ Controller = (*RemoveQueryController)(nil)

func (c *RemoveQueryController) Handle(_ context.Context, msg *nats.Msg) error {
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

	err = c.rulesPort.RemoveQueryRule(id)
	if errors.Is(err, types.ErrRuleNotExists) {
		_ = broker.RespondToMsg(msg, dto.RuleNotFound())
		return nil
	}
	if err != nil {
		_ = broker.RespondToMsg(msg, dto.InternalError())
		return nil
	}

	_ = broker.RespondToMsg(msg, "OK")
	return nil
}

func (c *RemoveQueryController) Subject() broker.Subject {
	return "notification.queries.remove"
}
