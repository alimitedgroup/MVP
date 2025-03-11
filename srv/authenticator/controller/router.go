package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib"
)

type ControllerRouter []lib.BrokerRoute

func NewControllerRouter(ar *authRouter) *ControllerRouter {
	return &ControllerRouter{
		ar,
	}
}

func (r ControllerRouter) Setup(ctx context.Context) error {
	for _, v := range r {
		err := v.Setup(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
