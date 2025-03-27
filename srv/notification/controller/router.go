package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib"
)

type ControllerRouter []lib.BrokerRoute

func NewControllerRouter(notificationRouter *notificationRouter) *ControllerRouter {
	return &ControllerRouter{
		notificationRouter,
	}
}

func (r ControllerRouter) Setup(ctx context.Context) error {
	for _, route := range r {
		if err := route.Setup(ctx); err != nil {
			return err
		}
	}
	return nil
}
