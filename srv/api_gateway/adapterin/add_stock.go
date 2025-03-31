package adapterin

import (
	"log/slog"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/gin-gonic/gin"
)

func NewAddStockController(warehouses portin.Warehouses) *AddStockController {
	return &AddStockController{warehouses: warehouses}
}

type AddStockController struct {
	warehouses portin.Warehouses
}

func (c *AddStockController) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.AddStockRequest
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

		err := c.warehouses.AddStock(req.WarehouseID, req.GoodID, req.Quantity)
		if err != nil {
			slog.Error("error while handling request to /api/v1/goods/:good_id/warehouse/:warehouse_id/stock", "error", err)
			ctx.JSON(500, dto.InternalError())
			return
		}
		ctx.JSON(200, map[string]interface{}{})
	}
}

func (c *AddStockController) Pattern() string {
	return "/goods/:good_id/warehouse/:warehouse_id/stock"
}

func (c *AddStockController) Method() string {
	return "POST"
}

func (c *AddStockController) RequiresAuth() bool {
	return true
}

func (c *AddStockController) AllowedRoles() []types.UserRole {
	return []types.UserRole{types.RoleLocalAdmin}
}

var _ Controller = (*AddStockController)(nil)
