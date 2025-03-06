package controller

import (
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portout"
	"github.com/gin-gonic/gin"
)

type LoginController struct {
	broker   *broker.NatsMessageBroker
	business *business.Business
}

func NewLoginController(broker *broker.NatsMessageBroker, business *business.Business) *LoginController {
	return &LoginController{broker, business}
}

func (c *LoginController) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.PostForm("username")
		token, err := c.business.Login(username)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if token.Token == "" {
			ctx.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		var role string
		switch token.Role {
		case portout.RoleClient:
			role = "client"
		case portout.RoleLocalAdmin:
			role = "local admin"
		case portout.RoleGlobalAdmin:
			role = "global admin"
		default:
			ctx.JSON(500, gin.H{"error": "unknown role"})
			return
		}

		ctx.JSON(200, dto.AuthLoginResponse{Token: string(token.Token), Role: role})
	}
}

func (c *LoginController) Pattern() string {
	return "/login"
}

func (c *LoginController) Method() string {
	return "POST"
}
