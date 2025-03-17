package sender

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
)

type PublishReservationEventAdapter struct {
	broker       *broker.NatsMessageBroker
	warehouseCfg *config.WarehouseConfig
}

func NewPublishReservationEventAdapter(broker *broker.NatsMessageBroker, warehouseCfg *config.WarehouseConfig) *PublishReservationEventAdapter {
	return &PublishReservationEventAdapter{broker, warehouseCfg}
}

func (a *PublishReservationEventAdapter) StoreReservationEvent(ctx context.Context, reservation model.Reservation) error {
	goods := make([]stream.ReservationGood, 0, len(reservation.Goods))
	for _, good := range reservation.Goods {
		goods = append(goods, stream.ReservationGood{
			GoodID:   string(good.GoodID),
			Quantity: good.Quantity,
		})
	}

	streamMsg := stream.ReservationEvent{
		Id:          string(reservation.ID),
		Goods:       goods,
		WarehouseID: a.warehouseCfg.ID,
	}

	payload, err := json.Marshal(streamMsg)
	if err != nil {
		return err
	}

	resp, err := a.broker.Js.Publish(ctx, fmt.Sprintf("reservation.%s", a.warehouseCfg.ID), payload)
	if err != nil {
		return err
	}

	_ = resp

	return nil
}
