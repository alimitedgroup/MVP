package listener

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
)

type CatalogRouter struct {
	catalogListener *CatalogListener
	broker          *broker.NatsMessageBroker
	restore         broker.IRestoreStreamControl
}

func NewCatalogRouter(catalogListener *CatalogListener, restoreFactory broker.IRestoreStreamControlFactory, broker *broker.NatsMessageBroker) *CatalogRouter {
	return &CatalogRouter{catalogListener, broker, restoreFactory.Build()}
}

func (r *CatalogRouter) Setup(ctx context.Context) error {
	// register stream message handlers
	err := r.broker.RegisterJsHandler(ctx, r.restore, stream.AddOrChangeGoodDataStream, r.catalogListener.ListenGoodUpdate)
	if err != nil {
		return err
	}

	// wait restoring of the state before starting the server
	r.restore.Wait()

	// register request/reply handlers

	return nil
}
