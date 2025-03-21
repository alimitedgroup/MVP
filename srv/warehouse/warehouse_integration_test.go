package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap/zaptest"
	"log"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/controller"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/listener"
	"github.com/alimitedgroup/MVP/srv/warehouse/application"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/magiconair/properties/assert"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
)

var Module = fx.Options(
	lib.Module,
	adapter.Module,
	application.Module,
)

type TestParams struct {
	Js jetstream.JetStream
	Ns *nats.Conn
}

func IntegrationTest(t *testing.T, testFunc any) {
	ctx := t.Context()

	ns, _ := broker.NewInProcessNATSServer(t)
	cfg := config.WarehouseConfig{ID: "1"}

	js, err := jetstream.New(ns)
	if err != nil {
		t.Error(err)
	}

	p := TestParams{
		Ns: ns,
		Js: js,
	}

	app := fx.New(
		Module,
		fx.Supply(&cfg),
		fx.Supply(&p),
		fx.Supply(ns),
		fx.Invoke(testFunc),
		fx.Supply(zaptest.NewLogger(t)),
	)
	err = app.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = app.Stop(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func NatsSendRequest[RequestDTO any, ResponseDTO any](t *testing.T, p *TestParams, subject string, dto RequestDTO) ResponseDTO {
	payload, err := json.Marshal(dto)
	if err != nil {
		t.Error(err)
	}

	msg, err := p.Ns.Request(subject, payload, 1*time.Second)
	if err != nil {
		t.Error(err)
	}

	var respDto ResponseDTO
	err = json.Unmarshal(msg.Data, &respDto)
	if err != nil {
		t.Error(err)
	}

	return respDto
}

func NatsPublishToStream[RequestDTO any](ctx context.Context, t *testing.T, p *TestParams, stream jetstream.StreamConfig, subject string, event RequestDTO) {
	s, err := p.Js.CreateStream(ctx, stream)
	if err != nil {
		t.Errorf("failed to create stream: %v", err)
	}

	beforeInfo, err := s.Info(ctx)
	if err != nil {
		t.Error(err)
	}

	payload, err := json.Marshal(event)
	if err != nil {
		t.Error(err)
	}

	ack, err := p.Js.Publish(ctx, subject, payload)
	if err != nil {
		t.Error(err)
	}

	info, err := s.Info(ctx)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, info.State.LastSeq, beforeInfo.State.LastSeq+1)
	assert.Equal(t, ack.Stream, stream.Name)
}

func TestAddAndRemoveWarehouseStock(t *testing.T) {
	IntegrationTest(t, func(
		lc fx.Lifecycle, cfg *config.WarehouseConfig, p *TestParams,
		stockControllerRouter *controller.StockRouter, catalogListenerRouter *listener.CatalogRouter, stockListenerRouter *listener.StockUpdateRouter) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				for _, v := range []lib.BrokerRoute{catalogListenerRouter, stockListenerRouter, stockControllerRouter} {
					if err := v.Setup(ctx); err != nil {
						t.Error(err)
					}
				}

				{
					event := stream.GoodUpdateData{
						GoodID:             "1",
						GoodNewName:        "hat",
						GoodNewDescription: "very good hat",
					}

					NatsPublishToStream(ctx, t, p, stream.AddOrChangeGoodDataStream, "good.update", event)
					time.Sleep(10 * time.Millisecond)
				}

				{
					dto := request.AddStockRequestDTO{GoodID: "1", Quantity: 10}
					resp := NatsSendRequest[request.AddStockRequestDTO, response.ResponseDTO[string]](
						t, p, fmt.Sprintf("warehouse.%s.stock.add", cfg.ID), dto,
					)

					assert.Equal(t, resp.Error, "")
					assert.Equal(t, resp.Message, "ok")
				}

				time.Sleep(10 * time.Millisecond)

				{
					dto := request.RemoveStockRequestDTO{GoodID: "1", Quantity: 5}
					resp := NatsSendRequest[request.RemoveStockRequestDTO, response.ResponseDTO[string]](
						t, p, fmt.Sprintf("warehouse.%s.stock.remove", cfg.ID), dto,
					)

					assert.Equal(t, resp.Error, "")
					assert.Equal(t, resp.Message, "ok")
				}

				time.Sleep(10 * time.Millisecond)

				{
					s, err := p.Js.Stream(ctx, stream.StockUpdateStreamConfig.Name)
					if err != nil {
						t.Error(err)
					}

					msg, err := s.GetLastMsgForSubject(ctx, fmt.Sprintf("stock.update.%s", cfg.ID))
					if err != nil {
						t.Error(err)
					}

					var event stream.StockUpdate
					if err = json.Unmarshal(msg.Data, &event); err != nil {
						t.Error(err)
					}

					assert.Equal(t, event.Goods[0].GoodID, "1")
					assert.Equal(t, event.Goods[0].Quantity, int64(5))
					assert.Equal(t, event.WarehouseID, cfg.ID)
					assert.Equal(t, event.Type, stream.StockUpdateTypeRemove)
				}

				return nil
			},
		})

	})
}
