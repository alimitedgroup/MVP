package model

type Reservation struct {
	ID    string
	Goods []ReservationGood
}

type ReservationID string

type ReservationGood struct {
	GoodID   string
	Quantity int64
}
