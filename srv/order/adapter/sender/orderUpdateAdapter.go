package sender

import (
	"context"
	"encoding/json"
	"time"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
)

type PublishOrderUpdateAdapter struct {
	broker *broker.NatsMessageBroker
}

func NewPublishOrderUpdateAdapter(broker *broker.NatsMessageBroker) *PublishOrderUpdateAdapter {
	return &PublishOrderUpdateAdapter{broker}
}

func (a *PublishOrderUpdateAdapter) SaveOrderUpdate(ctx context.Context, cmd port.SaveOrderUpdateCmd) error {
	now := time.Now()

	goods := make([]stream.OrderUpdateGood, 0)
	for _, good := range cmd.Goods {
		goods = append(goods, stream.OrderUpdateGood{
			GoodID:   good.GoodId,
			Quantity: good.Quantity,
		})
	}
	streamMsg := stream.OrderUpdate{
		ID:           cmd.ID,
		Name:         cmd.Name,
		Email:        cmd.Email,
		Address:      cmd.Address,
		Goods:        goods,
		CreationTime: now.UnixNano(),
		UpdateTime:   now.UnixNano(),
	}

	payload, err := json.Marshal(streamMsg)
	if err != nil {
		return err
	}

	resp, err := a.broker.Js.Publish(ctx, "order.update", payload)
	if err != nil {
		return err
	}

	_ = resp

	return nil
}
