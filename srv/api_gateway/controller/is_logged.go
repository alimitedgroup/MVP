package controller

import (
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/gin-gonic/gin"
)

type AuthHealthCheckController struct {
	business portin.Auth
}

func NewAuthHealthCheckController(business portin.Auth) *AuthHealthCheckController {
	return &AuthHealthCheckController{business: business}
}

func (c *AuthHealthCheckController) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, dto.IsLoggedResponse{
			Role: c.MustGet("user_data").(portin.UserData).Role.String(),
		})
	}
}

func (c *AuthHealthCheckController) Pattern() string {
	return "/is_logged"
}

func (c *AuthHealthCheckController) Method() string {
	return "GET"
}

func (c *AuthHealthCheckController) RequiresAuth() bool {
	return true
}

var _ Controller = (*AuthHealthCheckController)(nil)
