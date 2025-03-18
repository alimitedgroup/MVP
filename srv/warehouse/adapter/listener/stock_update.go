package listener

import (
	"context"
	"encoding/json"
	"log"

	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
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
	var cmdType port.StockUpdateCmdType
	switch event.Type {
	case stream.StockUpdateTypeAdd:
		cmdType = port.StockUpdateCmdTypeAdd
	case stream.StockUpdateTypeRemove:
		cmdType = port.StockUpdateCmdTypeRemove
	case stream.StockUpdateTypeOrder:
		cmdType = port.StockUpdateCmdTypeOrder
	case stream.StockUpdateTypeTransfer:
		cmdType = port.StockUpdateCmdTypeTransfer
	default:
		log.Fatal("unknown stock update type")
	}

	cmd := port.StockUpdateCmd{
		ID:            event.ID,
		Type:          cmdType,
		OrderID:       event.OrderID,
		TransferID:    event.TransferID,
		ReservationID: event.ReservationID,
		Timestamp:     event.Timestamp,
	}

	cmd.Goods = make([]port.StockUpdateCmdGood, 0, len(event.Goods))

	for _, good := range event.Goods {
		cmd.Goods = append(cmd.Goods, port.StockUpdateCmdGood{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
			Delta:    good.Delta,
		})
	}

	return cmd
}
