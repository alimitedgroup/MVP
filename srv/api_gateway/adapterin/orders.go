package adapterin

import (
	"log/slog"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/gin-gonic/gin"
)

func NewGetOrdersController(order portin.Order) *GetOrdersController {
	return &GetOrdersController{order: order}
}

type GetOrdersController struct {
	order portin.Order
}

func (c *GetOrdersController) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orders, err := c.order.GetOrders()
		if err != nil {
			slog.Error("error while handling request to /api/v1/orders", "error", err)
			ctx.JSON(500, dto.InternalError())
			return
		}
		ctx.JSON(200, dto.GetOrdersResponse{
			Orders: orders,
		})
	}
}

func (c *GetOrdersController) Pattern() string {
	return "/orders"
}

func (c *GetOrdersController) Method() string {
	return "GET"
}

func (c *GetOrdersController) RequiresAuth() bool {
	return false
}

var _ Controller = (*GetOrdersController)(nil)
