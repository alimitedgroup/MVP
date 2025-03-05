package controller

import (
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business"
	"github.com/gin-gonic/gin"
)

type NatsRouter struct {
	http     *lib.HTTPHandler
	broker   *broker.NatsMessageBroker
	business *business.Business
}

func NewNatsRouter(http *lib.HTTPHandler, broker *broker.NatsMessageBroker, business *business.Business) *NatsRouter {
	return &NatsRouter{http, broker, business}
}

func (nr *NatsRouter) Setup() {
	nr.http.ApiGroup.POST("/login", nr.Login)
}

func (nr *NatsRouter) Login(c *gin.Context) {
	username := c.PostForm("username")
	token, err := nr.business.Login(username)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if token == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}
