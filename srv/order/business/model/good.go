package model

type GoodId string

type GoodStock struct {
	ID       GoodId
	Quantity int64
}

type GoodInfo struct {
	ID          GoodId
	Name        string
	Description string
}
