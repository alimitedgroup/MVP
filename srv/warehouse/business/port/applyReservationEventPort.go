package port

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
)

type IApplyReservationEventPort interface {
	ApplyReservationEvent(model.Reservation) error
	ApplyOrderFilled(model.Reservation) error
}
