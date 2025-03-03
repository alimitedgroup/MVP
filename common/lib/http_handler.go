package lib

import "github.com/gin-gonic/gin"

type HTTPHandler struct {
	Engine   *gin.Engine
	ApiGroup *gin.RouterGroup
}

func NewHTTPHandler() *HTTPHandler {
	r := gin.New()
	api := r.Group("/api/v1")

	return &HTTPHandler{r, api}
}
