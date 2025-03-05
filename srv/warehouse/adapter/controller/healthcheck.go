package controller

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/nats-io/nats.go"
)

type HealthcheckController struct {
}

func NewHealthcheckController() *HealthcheckController {
	return &HealthcheckController{}
}

func (c *HealthcheckController) PingHandler(ctx context.Context, msg *nats.Msg) error {
	resp := response.HealthCheckResponseDTO{Message: "pong"}

	payload, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	err = msg.RespondMsg(&nats.Msg{Data: payload})
	if err != nil {
		return err
	}

	return nil
}
