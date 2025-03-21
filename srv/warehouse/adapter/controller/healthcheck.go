package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/nats-io/nats.go"
)

type HealthCheckController struct {
}

func NewHealthCheckController() *HealthCheckController {
	return &HealthCheckController{}
}

func (c *HealthCheckController) PingHandler(ctx context.Context, msg *nats.Msg) error {
	resp := response.HealthCheckResponseDTO{Message: "pong"}
	if err := broker.RespondToMsg(msg, resp); err != nil {
		return err
	}

	return nil
}
