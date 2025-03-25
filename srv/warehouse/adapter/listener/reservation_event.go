package listener

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/nats-io/nats.go/jetstream"
)

type ReservationEventListener struct {
	applyReservationEventUseCase port.IApplyReservationUseCase
}

func NewReservationEventListener(applyReservationEventUseCase port.IApplyReservationUseCase) *ReservationEventListener {
	return &ReservationEventListener{applyReservationEventUseCase}
}

func (l *ReservationEventListener) ListenReservationEvent(ctx context.Context, msg jetstream.Msg) error {
	var event stream.ReservationEvent
	err := json.Unmarshal(msg.Data(), &event)
	if err != nil {
		return err
	}

	goods := make([]port.ReservationGood, 0, len(event.Goods))
	for _, good := range event.Goods {
		goods = append(goods, port.ReservationGood(good))
	}
	cmd := port.ApplyReservationEventCmd{
		ID:    event.ID,
		Goods: goods,
	}
	err = l.applyReservationEventUseCase.ApplyReservationEvent(cmd)
	if err != nil {
		return err
	}

	return nil
}
