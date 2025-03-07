package listener

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

type goodMock struct {
	id          string
	name        string
	description string
}

type applyCatalogUpdateMock struct {
	sync.Mutex
	catalogMap map[string]goodMock
}

func newApplyCatalogUpdateMock() *applyCatalogUpdateMock {
	return &applyCatalogUpdateMock{catalogMap: make(map[string]goodMock)}
}

func (m *applyCatalogUpdateMock) ApplyCatalogUpdate(ctx context.Context, cmd port.CatalogUpdateCmd) error {
	m.Lock()
	defer m.Unlock()

	m.catalogMap[cmd.GoodId] = goodMock{
		id:          cmd.GoodId,
		name:        cmd.Name,
		description: cmd.Description,
	}

	return nil
}

func (m *applyCatalogUpdateMock) GetGood(id string) goodMock {
	m.Lock()
	defer m.Unlock()

	good := m.catalogMap[id]
	return good
}

func TestCatalogUpdateListener(t *testing.T) {
	ctx := t.Context()

	ns := broker.NewInProcessNATSServer(t)
	js, err := jetstream.New(ns)
	if err != nil {
		t.Error(err)
	}

	cfg := config.WarehouseConfig{ID: "1"}
	mock := newApplyCatalogUpdateMock()

	app := fx.New(
		fx.Supply(&cfg),
		fx.Supply(ns),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(fx.Annotate(broker.NewRestoreStreamControlFactory, fx.As(new(broker.IRestoreStreamControlFactory)))),
		fx.Supply(fx.Annotate(mock, fx.As(new(port.ApplyCatalogUpdateUseCase)))),
		fx.Provide(NewCatalogListener),
		fx.Provide(NewCatalogRouter),
		fx.Invoke(func(lc fx.Lifecycle, r *CatalogRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}

					event := stream.GoodUpdateData{
						GoodID:             "1",
						GoodNewName:        "hat",
						GoodNewDescription: "very nice hat",
					}

					payload, err := json.Marshal(event)
					if err != nil {
						t.Error(err)
					}

					ack, err := js.Publish(ctx, "good.update", payload)
					if err != nil {
						t.Error(err)
					}

					time.Sleep(100 * time.Millisecond)

					assert.Equal(t, ack.Stream, "good_data_update")
					assert.Equal(t, mock.GetGood("1").name, "hat")

					return nil
				},
			})
		}),
	)

	err = app.Start(ctx)
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
