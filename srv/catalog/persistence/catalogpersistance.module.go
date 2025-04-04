package persistence

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewCatalogRepository,
			fx.As(new(IGoodRepository)),
		),
		fx.Annotate(NewCatalogGoodDataRepository,
			fx.As(new(ICatalogGoodDataRepository)),
		),
	),
)
