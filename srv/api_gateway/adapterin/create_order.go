package adapterin

import (
	"log/slog"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/gin-gonic/gin"
)

func NewCreateOrderController(order portin.Order) *CreateOrderController {
	return &CreateOrderController{order: order}
}

type CreateOrderController struct {
	order portin.Order
}

func (c *CreateOrderController) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderId, err := c.order.CreateOrder("")
		if err != nil {
			slog.Error("error while handling request to /api/v1/order", "error", err)
			ctx.JSON(500, dto.InternalError())
			return
		}
		ctx.JSON(200, dto.CreateOrderResponse{
			OrderID: orderId,
		})
	}
}

func (c *CreateOrderController) Pattern() string {
	return "/order"
}

func (c *CreateOrderController) Method() string {
	return "POST"
}

func (c *CreateOrderController) RequiresAuth() bool {
	return false
}

var _ Controller = (*CreateOrderController)(nil)
