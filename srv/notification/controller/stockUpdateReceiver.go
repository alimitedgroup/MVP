package controller

import (
	"context"
	"encoding/json"
	"github.com/alimitedgroup/MVP/common/stream"
	servicecmd "github.com/alimitedgroup/MVP/srv/notification/service/cmd"
	serviceportin "github.com/alimitedgroup/MVP/srv/notification/service/portin"
	"github.com/nats-io/nats.go/jetstream"
	"time"
)

func NewStockUpdateReceiver(addStockUpdateUseCase serviceportin.IAddStockUpdateUseCase) *StockUpdateReceiver {
	return &StockUpdateReceiver{
		addStockUpdateUseCase: addStockUpdateUseCase,
	}
}

type StockUpdateReceiver struct {
	addStockUpdateUseCase serviceportin.IAddStockUpdateUseCase
}

var _ JsController = (*StockUpdateReceiver)(nil)

func (s StockUpdateReceiver) Handle(_ context.Context, msg jetstream.Msg) error {
	request := &stream.StockUpdate{}

	err := json.Unmarshal(msg.Data(), request)
	if err != nil {
		return err
	}

	goods := make([]servicecmd.StockGood, len(request.Goods))
	for i, g := range request.Goods {
		goods[i] = servicecmd.StockGood{
			ID:       g.GoodID,
			Quantity: int(g.Quantity),
			Delta:    int(g.Delta),
		}
	}

	cmd := servicecmd.NewAddStockUpdateCmd(request.WarehouseID, string(request.Type), request.OrderID, request.TransferID, goods, time.Now().Unix())
	_, err = s.addStockUpdateUseCase.AddStockUpdate(cmd)

	return err
}

func (s StockUpdateReceiver) Stream() jetstream.StreamConfig {
	return stream.StockUpdateStreamConfig
}
