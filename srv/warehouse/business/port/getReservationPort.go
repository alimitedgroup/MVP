package port

import "github.com/alimitedgroup/MVP/srv/warehouse/business/model"

type IGetReservationPort interface {
	GetReservation(reservationId model.ReservationID) (model.Reservation, error)
}
