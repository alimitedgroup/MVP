package adapterin

import (
	"context"
	"encoding/json"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/notification/portin"
	servicecmd "github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/nats-io/nats.go/jetstream"
	"time"
)

func NewStockUpdateReceiver(addStockUpdateUseCase portin.StockUpdates) *StockUpdateReceiver {
	return &StockUpdateReceiver{
		addStockUpdateUseCase: addStockUpdateUseCase,
	}
}

type StockUpdateReceiver struct {
	addStockUpdateUseCase portin.StockUpdates
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
	return s.addStockUpdateUseCase.RecordStockUpdate(cmd)
}

func (s StockUpdateReceiver) Stream() jetstream.StreamConfig {
	return stream.StockUpdateStreamConfig
}
