package port

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
)

type IStoreReservationEventPort interface {
	StoreReservationEvent(context.Context, model.Reservation) error
}
