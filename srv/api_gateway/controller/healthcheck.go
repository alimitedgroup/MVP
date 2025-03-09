package controller

import (
	"github.com/gin-gonic/gin"
)

type HealthCheckController struct {
}

func NewHealthCheckController() *HealthCheckController {
	return &HealthCheckController{}
}

func (c *HealthCheckController) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	}
}

func (c *HealthCheckController) Pattern() string {
	return "/ping"
}

func (c *HealthCheckController) Method() string {
	return "GET"
}

func (c *HealthCheckController) RequiresAuth() bool {
	return false
}

var _ Controller = (*HealthCheckController)(nil)
