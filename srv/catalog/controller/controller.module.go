package controller

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewCatalogGoodInfoController),
	fx.Provide(NewCatalogController),
	fx.Provide(NewCatalogRouter),
	fx.Provide(NewControllerRouter),
)
