package adapterin

import (
	"context"
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/notification/portin"
	"github.com/nats-io/nats.go"
)

func NewListQueriesController(rulesPort portin.QueryRules) *ListQueriesController {
	return &ListQueriesController{rulesPort: rulesPort}
}

type ListQueriesController struct {
	rulesPort portin.QueryRules
}

// Asserzione a compile-time che ListQueriesController implementi Controller
var _ Controller = (*ListQueriesController)(nil)

func (c *ListQueriesController) Handle(_ context.Context, msg *nats.Msg) error {
	rules, err := c.rulesPort.ListQueryRules()
	if err != nil {
		_ = broker.RespondToMsg(msg, dto.InternalError())
	} else {
		_ = broker.RespondToMsg(msg, rules)
	}

	return nil
}

func (c *ListQueriesController) Subject() broker.Subject {
	return "notification.queries.list"
}
