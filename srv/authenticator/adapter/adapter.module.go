package adapter

import (
	serviceportout "github.com/alimitedgroup/MVP/srv/authenticator/service/portOut"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewAuthAdapter,
			fx.As(new(serviceportout.ICheckKeyPairExistance)),
			fx.As(new(serviceportout.IGetPemPrivateKeyPort)),
			fx.As(new(serviceportout.IGetPemPublicKeyPort)),
			fx.As(new(serviceportout.IStorePemKeyPair)),
		),
		fx.Annotate(NewAuthPublisherAdapter,
			fx.As(new(serviceportout.IPublishPort)),
		)))
