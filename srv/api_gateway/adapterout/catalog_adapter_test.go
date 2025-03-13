package adapterout

import (
	"context"
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/catalog/catalogAdapter"
	"github.com/alimitedgroup/MVP/srv/catalog/controller"
	goodRepository "github.com/alimitedgroup/MVP/srv/catalog/persistence"
	"github.com/alimitedgroup/MVP/srv/catalog/service"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"testing"
)

func startCatalog(t *testing.T, nc *nats.Conn) {
	catalogSvc := fx.New(
		lib.Module,
		controller.Module,
		goodRepository.Module,
		catalogAdapter.Module,
		service.Module,
		fx.Supply(nc),
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

	brk, err := broker.NewNatsMessageBroker(nc)
	require.NoError(t, err)
	catalog := NewCatalogAdapter(brk)

	_, err = catalog.ListGoods()
	require.NoError(t, err)
}

func TestListStock(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)
	startCatalog(t, nc)

	brk, err := broker.NewNatsMessageBroker(nc)
	require.NoError(t, err)
	catalog := NewCatalogAdapter(brk)

	_, err = catalog.ListStock()
	require.NoError(t, err)
}

func TestListWarehouses(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)
	startCatalog(t, nc)

	brk, err := broker.NewNatsMessageBroker(nc)
	require.NoError(t, err)
	catalog := NewCatalogAdapter(brk)

	_, err = catalog.ListWarehouses()
	require.NoError(t, err)
}

func Run(cr *controller.ControllerRouter) error {
	return cr.Setup(context.Background())
}
