package adapterin

import (
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/gin-gonic/gin"
)

type LoginController struct {
	broker   *broker.NatsMessageBroker
	business portin.Auth
}

func NewLoginController(broker *broker.NatsMessageBroker, business portin.Auth) *LoginController {
	return &LoginController{broker, business}
}

func (c *LoginController) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, ok := ctx.GetPostForm("username")
		if !ok {
			ctx.JSON(400, dto.FieldIsRequired("username"))
			return
		}

		token, err := c.business.Login(username)
		if err != nil {
			ctx.JSON(500, dto.InternalError())
			return
		}

		if token.Token == "" {
			ctx.JSON(401, dto.AuthFailed())
			return
		}

		ctx.JSON(200, dto.AuthLoginResponse{Token: string(token.Token)})
	}
}

func (c *LoginController) Pattern() string {
	return "/login"
}

func (c *LoginController) Method() string {
	return "POST"
}

func (c *LoginController) RequiresAuth() bool {
	return false
}

func (c *LoginController) AllowedRoles() []types.UserRole {
	return []types.UserRole{types.RoleNone}
}

var _ Controller = (*LoginController)(nil)
