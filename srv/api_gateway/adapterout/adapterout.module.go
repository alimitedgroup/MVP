package adapterout

import "go.uber.org/fx"

var Module = fx.Module("adapterout", fx.Provide(NewAuthenticationAdapter))
