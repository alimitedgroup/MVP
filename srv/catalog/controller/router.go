package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib"
)

type ControllerRouter []lib.BrokerRoute

func NewControllerRouter(catalogRouter *catalogRouter) *ControllerRouter {
	return &ControllerRouter{
		catalogRouter,
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
