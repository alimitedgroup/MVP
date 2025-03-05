package controller

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type catalogController struct {
}

func NewCatalogController() *catalogController {
	return &catalogController{}
}

func (cc *catalogController) getGoodRequest(ctx context.Context, msg *nats.Msg) error {
	return nil
}

func (cc *catalogController) getWarehouseRequest(ctx context.Context, msg *nats.Msg) error {
	return nil
}

func (cc *catalogController) setGoodDataRequest(ctx context.Context, msg jetstream.Msg) error {
	return nil
}

func (cc *catalogController) setGoodQuantityRequest(ctx context.Context, msg jetstream.Msg) error {
	return nil
}
