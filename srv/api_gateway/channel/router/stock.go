package router

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/api_gateway/channel/constants"
	"github.com/alimitedgroup/MVP/srv/api_gateway/channel/controller"
)

type StockRouter struct {
	stockController *controller.StockController
	broker          *broker.NatsMessageBroker
}

func NewStockUpdateRouter(stockUpdateController *controller.StockController, n *broker.NatsMessageBroker) *StockRouter {
	return &StockRouter{stockUpdateController, n}
}

func (r *StockRouter) Setup(ctx context.Context) {
	_, _ = r.broker.RequestSubscribe(ctx, broker.StockUpdateSubject, constants.ApiGatewayGroup, r.stockController.UpdateHandler)
}
