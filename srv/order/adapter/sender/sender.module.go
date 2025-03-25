package sender

import (
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewNatsStreamAdapter,
		fx.As(new(port.ISendOrderUpdatePort)), fx.As(new(port.ISendContactWarehousePort)), fx.As(new(port.IRequestReservationPort)), fx.As(new(port.ISendTransferUpdatePort)),
	)),
)
