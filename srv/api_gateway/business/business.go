package business

import (
	"github.com/alimitedgroup/MVP/srv/api_gateway/portout"
	"go.uber.org/fx"
)

type Business struct {
	authAdapter portout.AuthenticationPortOut
}

func NewBusiness(authAdapter portout.AuthenticationPortOut) *Business {
	return &Business{authAdapter: authAdapter}
}

var Module = fx.Options(
	fx.Provide(NewBusiness),
)

func (b *Business) Login(username string) (portout.UserToken, error) {
	return b.authAdapter.GetToken(username)
}
