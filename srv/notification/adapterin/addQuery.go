package adapterin

import (
	"context"
	"encoding/json"
	"github.com/alimitedgroup/MVP/srv/notification/portin"
	servicecmd "github.com/alimitedgroup/MVP/srv/notification/types"
	"log"

	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/nats-io/nats.go/jetstream"
)

func NewAddQueryController(addQueryRuleUseCase portin.QueryRules) *AddQueryController {
	return &AddQueryController{addQueryRuleUseCase: addQueryRuleUseCase}
}

type AddQueryController struct {
	addQueryRuleUseCase portin.QueryRules
}

// Asserzione a compile-time che AddQueryController implementi JsController
var _ JsController = (*AddQueryController)(nil)

func (c *AddQueryController) Stream() jetstream.StreamConfig {
	return stream.QueryRuleStreamConfig
}

func (c *AddQueryController) Handle(_ context.Context, msg jetstream.Msg) error {
	log.Printf("addQueryRuleRequest ricevuto: %s", string(msg.Data()))

	request := stream.AddQueryRule{}
	err := json.Unmarshal(msg.Data(), &request)
	if err != nil {
		return err
	}

	cmd := servicecmd.QueryRule{GoodId: request.GoodID, Operator: request.Operator, Threshold: request.Threshold}
	response := c.addQueryRuleUseCase.AddQueryRule(&cmd)

	return response.GetOperationResult()
}
