package adapterin

import (
	"log/slog"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/gin-gonic/gin"
)

func NewUpdateGoodController(warehouses portin.Warehouses) *UpdateGoodController {
	return &UpdateGoodController{warehouses: warehouses}
}

type UpdateGoodController struct {
	warehouses portin.Warehouses
}

func (c *UpdateGoodController) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.UpdateGoodRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			slog.Error("invalid request body", "error", err)
			ctx.JSON(400, dto.InternalError())
			return
		}

		if err := ctx.ShouldBindUri(&req); err != nil {
			slog.Error("invalid request uri", "error", err)
			ctx.JSON(400, dto.InternalError())
			return
		}

		err := c.warehouses.UpdateGood(ctx, req.Id, req.Name, req.Description)
		if err != nil {
			slog.Error("error while handling request to /api/v1/goods", "error", err)
			ctx.JSON(500, dto.InternalError())
			return
		}
		ctx.JSON(200, nil)
	}
}

func (c *UpdateGoodController) Pattern() string {
	return "/goods/:good_id"
}

func (c *UpdateGoodController) Method() string {
	return "PUT"
}

func (c *UpdateGoodController) RequiresAuth() bool {
	return false
}

var _ Controller = (*UpdateGoodController)(nil)
