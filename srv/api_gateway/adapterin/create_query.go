package adapterin

import (
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewCreateQueryController(notifications portin.Notifications, logger *zap.Logger) *CreateQueryController {
	return &CreateQueryController{notifications: notifications, Logger: logger}
}

type CreateQueryController struct {
	notifications portin.Notifications
	*zap.Logger
}

func (c *CreateQueryController) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.CreateQueryRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			c.Error("invalid request body", zap.Error(err))
			ctx.JSON(400, dto.InternalError())
			return
		}

		queryId, err := c.notifications.CreateQuery(req.GoodID, req.Operator, req.Threshold)
		if err != nil {
			c.Error("error while creating notification query", zap.Error(err))
			ctx.JSON(500, dto.InternalError())
			return
		}
		ctx.JSON(200, dto.CreateQueryResponse{
			QueryID: queryId,
		})
	}
}

func (c *CreateQueryController) Pattern() string {
	return "/notifications/queries"
}

func (c *CreateQueryController) Method() string {
	return "POST"
}

func (c *CreateQueryController) RequiresAuth() bool {
	return true
}

var _ Controller = (*CreateQueryController)(nil)
