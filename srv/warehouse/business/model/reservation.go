package model

type Reservation struct {
	ID    ReservationId
	Goods []ReservationGood
}

type ReservationId string

type ReservationGood struct {
	GoodID   GoodID
	Quantity int64
}
