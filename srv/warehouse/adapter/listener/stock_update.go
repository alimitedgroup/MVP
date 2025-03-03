package listener

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"github.com/nats-io/nats.go/jetstream"
)

type StockUpdateListener struct {
	updateStockUseCase port.UpdateStockUseCase
}

func NewStockUpdateListener(updateStockUseCase port.UpdateStockUseCase) *StockUpdateListener {
	return &StockUpdateListener{updateStockUseCase}
}

func (l *StockUpdateListener) ListenStockUpdate(ctx context.Context, msg jetstream.Msg) error {
	var event stream.StockUpdate
	err := json.Unmarshal(msg.Data(), &event)
	if err != nil {
		return err
	}
	cmd := StockUpdateEventToUpdateStockCommand(event)
	err = l.updateStockUseCase.UpdateStock(cmd)
	if err != nil {
		return err
	}

	return nil
}

func StockUpdateEventToUpdateStockCommand(event stream.StockUpdate) port.UpdateStockCmd {
	cmd := port.UpdateStockCmd{
		ID:         event.ID,
		Type:       event.Type,
		OrderID:    event.OrderID,
		TransferID: event.TransferID,
		Timestamp:  event.Timestamp,
	}

	cmd.Goods = make([]port.UpdateStockCommandGood, 0, len(event.Goods))

	for _, good := range event.Goods {
		cmd.Goods = append(cmd.Goods, port.UpdateStockCommandGood{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
			Delta:    good.Delta,
		})
	}

	return cmd
}
