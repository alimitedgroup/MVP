package controller

import (
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"strings"
)

var Module = fx.Options(
	fx.Provide(AsController(NewHealthCheckController)),
	fx.Provide(AsController(NewLoginController)),
	fx.Invoke(fx.Annotate(RegisterRoutes, fx.ParamTags("", `group:"routes"`))),
)

func RegisterRoutes(http *lib.HTTPHandler, controllers []Controller) {
	for _, controller := range controllers {
		http.ApiGroup.Handle(controller.Method(), controller.Pattern(), controller.Handler())
	}
}

type Controller interface {
	Handler() gin.HandlerFunc
	Pattern() string
	Method() string
}

func AsController(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Controller)),
		fx.ResultTags(`group:"routes"`),
	)
}

func CheckRole(ctx *gin.Context, b *business.Business, roles []types.UserRole) {
	auth := ctx.GetHeader("Authorization")
	auth, found := strings.CutPrefix(auth, "Bearer ")
	if !found {
		ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthorized", "message": "No token provided"})
		return
	}
	data, err := b.ValidateToken(auth)
	if err != nil {
		ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthorized", "message": err.Error()})
		return
	}
	for _, role := range roles {
		if role == data.Role {
			return
		}
	}
	ctx.AbortWithStatusJSON(403, gin.H{"error": "forbidden", "message": "You don't have the required role"})
}
