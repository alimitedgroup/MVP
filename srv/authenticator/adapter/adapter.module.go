package adapter

import (
	serviceportout "github.com/alimitedgroup/MVP/srv/authenticator/service/portOut"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewAuthAdapter,
			fx.As(new(serviceportout.CheckKeyPairExistance)),
			fx.As(new(serviceportout.GetPemPrivateKeyPort)),
			fx.As(new(serviceportout.GetPemPublicKeyPort)),
			fx.As(new(serviceportout.IStorePemKeyPairInterface)),
		)))
