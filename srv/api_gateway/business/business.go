package business

import "go.uber.org/fx"

type Business struct {
}

func NewBusiness() *Business {
	return &Business{}
}

var Module = fx.Options(
	fx.Provide(NewBusiness),
)

func (*Business) Login(username string) string {

}
