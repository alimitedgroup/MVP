package port

type IApplyReservationUseCase interface {
	ApplyReservationEvent(ApplyReservationEventCmd) error
}

type ApplyReservationEventCmd struct {
	Id    string
	Goods []ReservationGood
}
