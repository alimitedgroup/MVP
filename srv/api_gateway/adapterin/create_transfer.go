package adapterin

import (
	"log/slog"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/gin-gonic/gin"
)

func NewCreateTransferController(order portin.Order) *CreateTransferController {
	return &CreateTransferController{order: order}
}

type CreateTransferController struct {
	order portin.Order
}

func (c *CreateTransferController) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		transferId, err := c.order.CreateTransfer("")
		if err != nil {
			slog.Error("error while handling request to /api/v1/transfer", "error", err)
			ctx.JSON(500, dto.InternalError())
			return
		}
		ctx.JSON(200, dto.CreateTransferResponse{
			TransferID: transferId,
		})
	}
}

func (c *CreateTransferController) Pattern() string {
	return "/transfer"
}

func (c *CreateTransferController) Method() string {
	return "POST"
}

func (c *CreateTransferController) RequiresAuth() bool {
	return false
}

var _ Controller = (*CreateTransferController)(nil)
