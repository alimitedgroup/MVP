package controller

import (
	"context"
	"encoding/json"
	"log"

	"github.com/alimitedgroup/MVP/common/stream"
	servicecmd "github.com/alimitedgroup/MVP/srv/notification/service/cmd"
	serviceportin "github.com/alimitedgroup/MVP/srv/notification/service/portin"
	"github.com/nats-io/nats.go/jetstream"
)

func NewAddQueryController(addQueryRuleUseCase serviceportin.IAddQueryRuleUseCase) *AddQueryController {
	return &AddQueryController{addQueryRuleUseCase: addQueryRuleUseCase}
}

type AddQueryController struct {
	addQueryRuleUseCase serviceportin.IAddQueryRuleUseCase
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

	cmd := servicecmd.NewAddQueryRuleCmd(request.GoodID, request.Operator, request.Threshold)
	response := c.addQueryRuleUseCase.AddQueryRule(cmd)

	return response.GetOperationResult()
}
