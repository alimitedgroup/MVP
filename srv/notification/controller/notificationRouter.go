package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
)

type notificationRouter struct {
	mb         *broker.NatsMessageBroker
	controller *notificationController
	rsc        *broker.RestoreStreamControl
}

func NewNotificationRouter(mb *broker.NatsMessageBroker, nc *notificationController, rsc *broker.RestoreStreamControl) *notificationRouter {
	return &notificationRouter{mb, nc, rsc}
}

func (nr *notificationRouter) Setup(ctx context.Context) error {
	err := nr.mb.RegisterJsHandler(ctx, nr.rsc, stream.StockUpdateStreamConfig, nr.controller.addStockUpdateRequest)
	if err != nil {
		return err
	}

	nr.rsc.Wait()

	err = nr.mb.RegisterRequest(ctx, "notification.addQueryRule", "notification", nr.controller.addQueryRuleRequest)
	if err != nil {
		return err
	}

	return nil
}
