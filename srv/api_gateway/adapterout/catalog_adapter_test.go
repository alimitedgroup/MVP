package adapterout

import (
	"context"
	"testing"

	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/catalog/catalogAdapter"
	"github.com/alimitedgroup/MVP/srv/catalog/controller"
	goodRepository "github.com/alimitedgroup/MVP/srv/catalog/persistence"
	"github.com/alimitedgroup/MVP/srv/catalog/service"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func startCatalog(t *testing.T, nc *nats.Conn) {
	logger := observability.TestLogger(t)
	catalogSvc := fx.New(
		lib.ModuleTest,
		controller.Module,
		goodRepository.Module,
		catalogAdapter.Module,
		service.Module,
		fx.Provide(observability.TestMeter),
		fx.Supply(nc, t, logger),
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
