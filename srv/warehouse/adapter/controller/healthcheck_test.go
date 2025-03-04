package controller_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/controller"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/magiconair/properties/assert"
	"go.uber.org/fx"
)

func TestHealthCheckController(t *testing.T) {
	ctx := t.Context()

	ns := broker.NewInProcessNATSServer(t)
	cfg := config.WarehouseConfig{ID: "1"}

	app := fx.New(
		fx.Supply(&cfg),
		fx.Supply(ns),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(controller.NewHealthcheckController),
		fx.Provide(controller.NewHealthCheckRouter),
		fx.Invoke(func(lc fx.Lifecycle, r *controller.HealthCheckRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}

					resp, err := ns.Request(fmt.Sprintf("warehouse.ping.%s", cfg.ID), []byte{}, 1*time.Second)
					if err != nil {
						t.Error(err)
					}

					var respDto response.HealthCheckResponseDTO
					err = json.Unmarshal(resp.Data, &respDto)
					if err != nil {
						t.Error(err)
					}

					assert.Equal(t, respDto.Message, "pong")

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
