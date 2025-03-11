package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib/broker"
)

type authRouter struct {
	mb         *broker.NatsMessageBroker
	controller *authController
	rsc        *broker.RestoreStreamControl
}

func NewAuthRouter(mb *broker.NatsMessageBroker, cc *authController, rsc *broker.RestoreStreamControl) *authRouter {
	return &authRouter{mb, cc, rsc}
}

func (ar *authRouter) Setup(ctx context.Context) error {
	/*var test []byte
	test = append(test, 7)
	ar.mb.Js.Publish(ctx, "key.ciao", test)*/
	err := ar.mb.RegisterRequest(ctx, "auth.login", "login", ar.controller.NewTokenRequest) //GetGoodsInfo
	if err != nil {
		return nil
	}
	return nil
}
