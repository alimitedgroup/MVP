package port

type IApplyReservationUseCase interface {
	ApplyReservationEvent(ApplyReservationEventCmd) error
}

type ApplyReservationEventCmd struct {
	ID    string
	Goods []ReservationGood
}
