package controller

import (
	"github.com/alimitedgroup/MVP/srv/api_gateway/business"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/gin-gonic/gin"
)

type AuthHealthCheckController struct {
	businness *business.Business
}

func NewAuthHealthCheckController(businness *business.Business) *AuthHealthCheckController {
	return &AuthHealthCheckController{businness: businness}
}

func (c *AuthHealthCheckController) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"role": c.MustGet("user_data").(portin.UserData).Role.String()})
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
