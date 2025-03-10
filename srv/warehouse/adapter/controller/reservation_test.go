package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/magiconair/properties/assert"
	"go.uber.org/fx"
)

func TestReservationController(t *testing.T) {
	t.Skip("not implemented")
	ctx := t.Context()

	ns := broker.NewInProcessNATSServer(t)
	cfg := config.WarehouseConfig{ID: "1"}

	app := fx.New(
		fx.Supply(&cfg),
		fx.Supply(ns),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(NewReservationController),
		fx.Provide(NewReservationRouter),
		fx.Invoke(func(lc fx.Lifecycle, r *ReservationRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}

					resp, err := ns.Request(fmt.Sprintf("warehouse.%s.reservation.add", cfg.ID), []byte{}, 1*time.Second)
					if err != nil {
						t.Error(err)
					}

					var respDto response.ResponseDTO[string]
					err = json.Unmarshal(resp.Data, &respDto)
					if err != nil {
						t.Error(err)
					}

					assert.Equal(t, respDto.Message, "ok")

					return nil
				},
				OnStop: func(ctx context.Context) error {
					return nil
				},
			})
		}),
	)

	err := app.Start(ctx)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = app.Stop(ctx)
		if err != nil {
			t.Error(err)
		}
	}()

}
