package controller

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/nats-io/nats.go"
)

type HealthCheckController struct {
}

func NewHealthCheckController() *HealthCheckController {
	return &HealthCheckController{}
}

func (c *HealthCheckController) PingHandler(ctx context.Context, msg *nats.Msg) error {
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
