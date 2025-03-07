package persistance

import (
	"github.com/alimitedgroup/MVP/srv/catalog/catalogAdapter"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewCatalogRepository,
		fx.As(new(catalogAdapter.IGoodRepository)))),
)
