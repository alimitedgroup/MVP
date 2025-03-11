package controller

import (
	"context"

	"github.com/nats-io/nats.go"
)

type authController struct {
}

func NewAuthController() *authController {
	return &authController{}
}

func (ar *authController) NewTokenRequest(ctx context.Context, msg *nats.Msg) error {
	return nil
}
