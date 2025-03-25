package listener

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/nats-io/nats.go/jetstream"
)

type StockUpdateListener struct {
	applyStockUpdateUseCase port.IApplyStockUpdateUseCase
}

func NewStockUpdateListener(applyStockUpdateUseCase port.IApplyStockUpdateUseCase) *StockUpdateListener {
	return &StockUpdateListener{applyStockUpdateUseCase}
}

func (l *StockUpdateListener) ListenStockUpdate(ctx context.Context, msg jetstream.Msg) error {
	var event stream.StockUpdate

	err := json.Unmarshal(msg.Data(), &event)
	if err != nil {
		return err
	}

	cmd := stockUpdateEventToApplyStockUpdateCmd(event)
	err = l.applyStockUpdateUseCase.ApplyStockUpdate(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}

func stockUpdateEventToApplyStockUpdateCmd(event stream.StockUpdate) port.StockUpdateCmd {
	goods := make([]port.StockUpdateGood, 0, len(event.Goods))

	for _, good := range event.Goods {
		goods = append(goods, port.StockUpdateGood{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
			Delta:    good.Delta,
		})
	}

	return port.StockUpdateCmd{
		ID:            event.ID,
		WarehouseID:   event.WarehouseID,
		Type:          port.StockUpdateType(event.Type),
		OrderID:       event.OrderID,
		TransferID:    event.TransferID,
		ReservationID: event.ReservationID,
		Timestamp:     event.Timestamp,
		Goods:         goods,
	}
}
