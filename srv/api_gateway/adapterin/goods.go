package adapterin

import (
	"log/slog"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/gin-gonic/gin"
)

func NewGetGoodsController(warehouses portin.Warehouses) *GetGoodsController {
	return &GetGoodsController{warehouses: warehouses}
}

type GetGoodsController struct {
	warehouses portin.Warehouses
}

func (c *GetGoodsController) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		goods, err := c.warehouses.GetGoods()
		if err != nil {
			slog.Error("error while handling request to /api/v1/goods", "error", err)
			ctx.JSON(500, dto.InternalError())
			return
		}
		ctx.JSON(200, dto.GetGoodsResponse{
			Goods: goods,
		})
	}
}

func (c *GetGoodsController) Pattern() string {
	return "/goods"
}

func (c *GetGoodsController) Method() string {
	return "GET"
}

func (c *GetGoodsController) RequiresAuth() bool {
	return true
}

func (c *GetGoodsController) AllowedRoles() []types.UserRole {
	return []types.UserRole{types.RoleNone}
}

var _ Controller = (*GetGoodsController)(nil)
