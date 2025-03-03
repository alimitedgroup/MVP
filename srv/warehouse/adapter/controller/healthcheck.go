package controller

import (
	"context"

	"github.com/nats-io/nats.go"
)

type HealthcheckController struct {
}

func NewHealthcheckController() *HealthcheckController {
	return &HealthcheckController{}
}

func (c *HealthcheckController) PingHandler(ctx context.Context, msg *nats.Msg) error {
	return nil
}
