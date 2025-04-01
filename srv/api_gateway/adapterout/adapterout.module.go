package adapterout

import (
	"github.com/alimitedgroup/MVP/common/lib/observability"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"adapterout",
	fx.Decorate(observability.WrapLogger("adapterout")),
	fx.Provide(NewAuthenticationAdapter),
	fx.Provide(NewCatalogAdapter),
	fx.Provide(NewOrderAdapter),
	fx.Provide(NewNotificationsAdapter),
)
