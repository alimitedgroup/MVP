package adapterin

import (
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewGetQueriesController(notifications portin.Notifications, logger *zap.Logger) *GetQueriesController {
	return &GetQueriesController{notifications: notifications, Logger: logger}
}

type GetQueriesController struct {
	notifications portin.Notifications
	*zap.Logger
}

func (c *GetQueriesController) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		queries, err := c.notifications.GetQueries()
		if err != nil {
			c.Error("Error while handling request to /api/v1/notifications/queries", zap.Error(err))
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
	return true
}

func (c *GetQueriesController) AllowedRoles() []types.UserRole {
	return []types.UserRole{types.RoleGlobalAdmin}
}

var _ Controller = (*GetQueriesController)(nil)
