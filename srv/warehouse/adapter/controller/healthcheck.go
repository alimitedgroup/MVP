package controller

import (
	"github.com/gin-gonic/gin"
)

type HealthcheckController struct {
}

func NewHealthcheckController() *HealthcheckController {
	return &HealthcheckController{}
}

func (c *HealthcheckController) PingHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"data": "pong",
	})
}
