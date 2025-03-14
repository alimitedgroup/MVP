package listener

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/nats-io/nats.go/jetstream"
)

type OrderUpdateListener struct {
	applyOrderUpdateUseCase port.IApplyOrderUpdateUseCase
}

func NewOrderUpdateListener(applyOrderUpdateUseCase port.IApplyOrderUpdateUseCase) *OrderUpdateListener {
	return &OrderUpdateListener{applyOrderUpdateUseCase}
}

func (l *OrderUpdateListener) ListenOrderUpdate(ctx context.Context, msg jetstream.Msg) error {
	var event stream.OrderUpdate
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		return err
	}

	cmd := orderUpdateEventToApplyOrderUpdateCmd(event)
	if err := l.applyOrderUpdateUseCase.ApplyOrderUpdate(ctx, cmd); err != nil {
		return err
	}

	return nil
}

func orderUpdateEventToApplyOrderUpdateCmd(event stream.OrderUpdate) port.OrderUpdateCmd {
	goods := make([]port.OrderUpdateGood, 0, len(event.Goods))
	for _, good := range event.Goods {
		goods = append(goods, port.OrderUpdateGood{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}

	cmd := port.OrderUpdateCmd{
		ID:           event.ID,
		Goods:        goods,
		Status:       event.Status,
		Name:         event.Name,
		Email:        event.Email,
		Address:      event.Address,
		CreationTime: event.CreationTime,
	}

	return cmd
}
