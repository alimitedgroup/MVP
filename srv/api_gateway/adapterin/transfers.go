package adapterin

import (
	"log/slog"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/gin-gonic/gin"
)

func NewGetTransfersController(order portin.Order) *GetTransfersController {
	return &GetTransfersController{order: order}
}

type GetTransfersController struct {
	order portin.Order
}

func (c *GetTransfersController) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		transfers, err := c.order.GetTransfers()
		if err != nil {
			slog.Error("error while handling request to /api/v1/transfers", "error", err)
			ctx.JSON(500, dto.InternalError())
			return
		}
		ctx.JSON(200, dto.GetTransfersResponse{
			Transfers: transfers,
		})
	}
}

func (c *GetTransfersController) Pattern() string {
	return "/transfers"
}

func (c *GetTransfersController) Method() string {
	return "GET"
}

func (c *GetTransfersController) RequiresAuth() bool {
	return false
}

var _ Controller = (*GetTransfersController)(nil)
