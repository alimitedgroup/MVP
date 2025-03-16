package listener

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/nats-io/nats.go/jetstream"
)

type OrderUpdateListener struct {
	confirmOrderUseCase port.IConfirmOrderUseCase
}

func NewOrderUpdateListener(confirmOrderUseCase port.IConfirmOrderUseCase) *OrderUpdateListener {
	return &OrderUpdateListener{confirmOrderUseCase}
}

func (l *OrderUpdateListener) ListenOrderUpdate(ctx context.Context, msg jetstream.Msg) error {
	var event stream.OrderUpdate
	err := json.Unmarshal(msg.Data(), &event)
	if err != nil {
		return err
	}

	goods := make([]port.OrderUpdateGood, 0, len(event.Goods))
	for _, good := range event.Goods {
		goods = append(goods, port.OrderUpdateGood{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}
	cmd := port.ConfirmOrderCmd{
		OrderID:      event.ID,
		Status:       event.Status,
		Reservations: event.Reservations,
		Goods:        goods,
	}
	err = l.confirmOrderUseCase.ConfirmOrder(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}
