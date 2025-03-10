package controller

import (
	"github.com/alimitedgroup/MVP/srv/api_gateway/business"
	"github.com/gin-gonic/gin"
	"strings"
)

type HTTPHandler struct {
	Engine             *gin.Engine
	ApiGroup           *gin.RouterGroup
	AuthenticatedGroup *gin.RouterGroup
}

func NewHTTPHandler(b *business.Business) *HTTPHandler {
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.Use(gin.Recovery())
	api := r.Group("/api/v1")
	authenticated := r.Group("/api/v1")
	authenticated.Use(Authentication(b))

	return &HTTPHandler{r, api, authenticated}
}

func Authentication(b *business.Business) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth, found := strings.CutPrefix(ctx.GetHeader("Authorization"), "Bearer ")
		if !found {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthorized", "message": "No token provided"})
			return
		}
		data, err := b.ValidateToken(auth)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthorized", "message": err.Error()})
			return
		}

		ctx.Set("user_data", data)
		ctx.Next()
	}
}
