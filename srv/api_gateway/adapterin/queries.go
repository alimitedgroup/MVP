package adapterin

import (
	"log/slog"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/gin-gonic/gin"
)

func NewGetQueriesController(notifications portin.Notifications) *GetQueriesController {
	return &GetQueriesController{notifications: notifications}
}

type GetQueriesController struct {
	notifications portin.Notifications
}

func (c *GetQueriesController) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		queries, err := c.notifications.GetQueries()
		if err != nil {
			slog.Error("error while handling request to /api/v1/notifications/queries", "error", err)
			ctx.JSON(500, dto.InternalError())
			return
		}
		ctx.JSON(200, dto.GetQueriesResponse{
			Queries: queries,
		})
	}
}

func (c *GetQueriesController) Pattern() string {
	return "/notifications/queries"
}

func (c *GetQueriesController) Method() string {
	return "GET"
}

func (c *GetQueriesController) RequiresAuth() bool {
	return false
}

var _ Controller = (*GetQueriesController)(nil)
