package adapterin

import (
	"log/slog"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/gin-gonic/gin"
)

func NewRemoveStockController(warehouses portin.Warehouses) *RemoveStockController {
	return &RemoveStockController{warehouses: warehouses}
}

type RemoveStockController struct {
	warehouses portin.Warehouses
}

func (c *RemoveStockController) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.RemoveStockRequest
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

		err := c.warehouses.RemoveStock(req.WarehouseID, req.GoodID, req.Quantity)
		if err != nil {
			slog.Error("error while handling request to /api/v1/goods/:good_id/warehouse/:warehouse_id/stock", "error", err)
			ctx.JSON(500, dto.InternalError())
			return
		}
		ctx.JSON(200, map[string]interface{}{})
	}
}

func (c *RemoveStockController) Pattern() string {
	return "/goods/:good_id/warehouse/:warehouse_id/stock"
}

func (c *RemoveStockController) Method() string {
	return "DELETE"
}

func (c *RemoveStockController) RequiresAuth() bool {
	return true
}

func (c *RemoveStockController) AllowedRoles() []types.UserRole {
	return []types.UserRole{types.RoleLocalAdmin}
}

var _ Controller = (*RemoveStockController)(nil)
