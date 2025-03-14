package listener

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewStockUpdateListener),
	fx.Provide(NewStockUpdateRouter),
	fx.Provide(NewListenerRoutes),
	fx.Provide(NewCatalogListener),
	fx.Provide(NewCatalogRouter),
	fx.Provide(NewReservationEventListener),
	fx.Provide(NewReservationEventRouter),
)
