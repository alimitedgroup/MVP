package adapterin

import (
	"log/slog"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/gin-gonic/gin"
)

func NewCreateQueryController(notifications portin.Notifications) *CreateQueryController {
	return &CreateQueryController{notifications: notifications}
}

type CreateQueryController struct {
	notifications portin.Notifications
}

func (c *CreateQueryController) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.CreateQueryRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			slog.Error("invalid request body", "error", err)
			ctx.JSON(400, dto.InternalError())
			return
		}

		queryId, err := c.notifications.CreateQuery(req.GoodID, req.Operator, req.Threshold)
		if err != nil {
			slog.Error("error while handling request to /api/v1/notifications/queries", "error", err)
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
