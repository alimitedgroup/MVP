package adapterin

import (
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/gin-gonic/gin"
	"log/slog"
)

func NewListWarehousesController(warehouses portin.Warehouses) *ListWarehousesController {
	return &ListWarehousesController{warehouses: warehouses}
}

type ListWarehousesController struct {
	warehouses portin.Warehouses
}

func Map[T, U any](seq []T, f func(T) U) []U {
	res := make([]U, 0, len(seq))
	for _, t := range seq {
		res = append(res, f(t))
	}
	return res
}

func (c *ListWarehousesController) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		warehouses, err := c.warehouses.GetWarehouses()
		if err != nil {
			slog.Error("error while handling request to /api/v1/warehouses", "error", err)
			ctx.JSON(500, dto.InternalError())
			return
		}
		ctx.JSON(200, dto.GetWarehousesResponse{
			Ids: Map(
				warehouses,
				func(w portin.WarehouseOverview) string { return w.ID },
			),
		})
	}
}

func (c *ListWarehousesController) Pattern() string {
	return "/warehouses"
}

func (c *ListWarehousesController) Method() string {
	return "GET"
}

func (c *ListWarehousesController) RequiresAuth() bool {
	return false
}

var _ Controller = (*ListWarehousesController)(nil)
