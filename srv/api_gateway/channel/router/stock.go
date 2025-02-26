package router

import (
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/api_gateway/channel/constants"
	"github.com/alimitedgroup/MVP/srv/api_gateway/channel/controller"
)

type StockUpdateRouter struct {
	stockController *controller.StockController
	broker          *broker.NatsMessageBroker
}

func NewStockUpdateRouter(stockUpdateController *controller.StockController, n *broker.NatsMessageBroker) *StockUpdateRouter {
	return &StockUpdateRouter{stockUpdateController, n}
}

func (r *StockUpdateRouter) Setup() {
	r.broker.RequestSubscribe(broker.StockUpdateSubject, constants.ApiGatewayGroup, r.stockController.UpdateHandler)
}
