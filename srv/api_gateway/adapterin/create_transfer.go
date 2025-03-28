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
		var req dto.CreateTransferRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			slog.Error("invalid request body", "error", err)
			ctx.JSON(400, dto.InternalError())
			return
		}

		transferId, err := c.order.CreateTransfer(req.SenderID, req.ReceiverID, req.Goods)
		if err != nil {
			slog.Error("error while handling request to /api/v1/transfers", "error", err)
			ctx.JSON(500, dto.InternalError())
			return
		}
		ctx.JSON(200, dto.CreateTransferResponse{
			TransferID: transferId,
		})
	}
}

func (c *CreateTransferController) Pattern() string {
	return "/transfers"
}

func (c *CreateTransferController) Method() string {
	return "POST"
}

func (c *CreateTransferController) RequiresAuth() bool {
	return true
}

var _ Controller = (*CreateTransferController)(nil)
