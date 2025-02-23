package controller

import (
	"github.com/gin-gonic/gin"
)

type HealthCheckController struct {
	// http *lib.HTTPHandler
}

func NewHealthCheckController() *HealthCheckController {
	return &HealthCheckController{}
}

func (c *HealthCheckController) Ping(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}
