package port

import "github.com/alimitedgroup/MVP/srv/warehouse/business/model"

type IGetReservationPort interface {
	GetReservation(reservationId model.ReservationId) (model.Reservation, error)
}
