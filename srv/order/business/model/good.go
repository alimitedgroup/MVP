package model

type GoodID string

type GoodStock struct {
	ID       GoodID
	Quantity int64
}

type GoodInfo struct {
	ID          GoodID
	Name        string
	Description string
}
