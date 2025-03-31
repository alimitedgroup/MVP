package adapterin

import (
	"log/slog"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/gin-gonic/gin"
)

func NewCreateGoodController(warehouses portin.Warehouses) *CreateGoodController {
	return &CreateGoodController{warehouses: warehouses}
}

type CreateGoodController struct {
	warehouses portin.Warehouses
}

func (c *CreateGoodController) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.CreateGoodRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			slog.Error("invalid request body", "error", err)
			ctx.JSON(400, dto.InternalError())
			return
		}

		goodId, err := c.warehouses.CreateGood(ctx, req.Name, req.Description)
		if err != nil {
			slog.Error("error while handling request to /api/v1/goods", "error", err)
			ctx.JSON(500, dto.InternalError())
			return
		}
		ctx.JSON(200, dto.CreateGoodResponse{
			GoodID: goodId,
		})
	}
}

func (c *CreateGoodController) Pattern() string {
	return "/goods"
}

func (c *CreateGoodController) Method() string {
	return "POST"
}

func (c *CreateGoodController) RequiresAuth() bool {
	return true
}

func (c *CreateGoodController) AllowedRoles() []types.UserRole {
	return []types.UserRole{types.RoleGlobalAdmin}
}

var _ Controller = (*CreateGoodController)(nil)
