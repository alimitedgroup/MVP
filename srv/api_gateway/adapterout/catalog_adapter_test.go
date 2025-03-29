package adapterout

import (
	"context"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/catalog/catalogAdapter"
	"github.com/alimitedgroup/MVP/srv/catalog/controller"
	goodRepository "github.com/alimitedgroup/MVP/srv/catalog/persistence"
	"github.com/alimitedgroup/MVP/srv/catalog/service"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func startCatalog(t *testing.T, nc *nats.Conn) {
	catalogSvc := fx.New(
		lib.ModuleTest,
		controller.Module,
		goodRepository.Module,
		catalogAdapter.Module,
		service.Module,
		fx.Supply(nc, t),
		fx.Invoke(Run),
	)

	err := catalogSvc.Start(context.Background())
	require.NoError(t, err)

	t.Cleanup(func() {
		err := catalogSvc.Stop(context.Background())
		require.NoError(t, err)
	})
}

func TestListGoods(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)
	startCatalog(t, nc)

	brk := broker.NewTest(t, nc)
	catalog := NewCatalogAdapter(brk)

	_, err := catalog.ListGoods()
	require.NoError(t, err)
}

func TestListStock(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)
	startCatalog(t, nc)

	brk := broker.NewTest(t, nc)
	catalog := NewCatalogAdapter(brk)

	_, err := catalog.ListStock()
	require.NoError(t, err)
}

func TestListWarehouses(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)
	startCatalog(t, nc)

	brk := broker.NewTest(t, nc)
	catalog := NewCatalogAdapter(brk)

	_, err := catalog.ListWarehouses()
	require.NoError(t, err)
}

func Run(cr *controller.ControllerRouter) error {
	return cr.Setup(context.Background())
}

func TestAddStock(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)

	sub, err := nc.Subscribe("warehouse.1.stock.add", func(msg *nats.Msg) {
		err := msg.Respond([]byte(`{"error": "", "message": "ok"}`))
		require.NoError(t, err)
	})
	require.NoError(t, err)
	defer func() {
		err := sub.Unsubscribe()
		require.NoError(t, err)
	}()

	brk := broker.NewTest(t, nc)
	catalog := NewCatalogAdapter(brk)

	err = catalog.AddStock("1", "1", 1)
	require.NoError(t, err)
}

func TestRemoveStock(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)

	sub, err := nc.Subscribe("warehouse.1.stock.remove", func(msg *nats.Msg) {
		err := msg.Respond([]byte(`{"error": "", "message": "ok"}`))
		require.NoError(t, err)
	})
	require.NoError(t, err)
	defer func() {
		err := sub.Unsubscribe()
		require.NoError(t, err)
	}()

	brk := broker.NewTest(t, nc)
	catalog := NewCatalogAdapter(brk)

	err = catalog.RemoveStock("1", "1", 1)
	require.NoError(t, err)
}

func TestCreateGood(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)

	js, err := jetstream.New(nc)
	require.NoError(t, err)

	s, err := js.CreateStream(t.Context(), stream.AddOrChangeGoodDataStream)
	require.NoError(t, err)
	require.Equal(t, "good_data_update", s.CachedInfo().Config.Name)

	brk := broker.NewTest(t, nc)
	catalog := NewCatalogAdapter(brk)

	goodId, err := catalog.CreateGood(t.Context(), "name", "description")
	require.NoError(t, err)
	require.NotEmpty(t, goodId)

	time.Sleep(100 * time.Millisecond)

	info, err := s.Info(t.Context())
	require.NoError(t, err)
	require.Equal(t, uint64(1), info.State.Msgs)
}

func TestUpdateGood(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)

	js, err := jetstream.New(nc)
	require.NoError(t, err)

	s, err := js.CreateStream(t.Context(), stream.AddOrChangeGoodDataStream)
	require.NoError(t, err)
	require.Equal(t, "good_data_update", s.CachedInfo().Config.Name)

	brk := broker.NewTest(t, nc)
	catalog := NewCatalogAdapter(brk)

	err = catalog.UpdateGood(t.Context(), "1", "name", "description")
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	info, err := s.Info(t.Context())
	require.NoError(t, err)
	require.Equal(t, uint64(1), info.State.Msgs)
}
