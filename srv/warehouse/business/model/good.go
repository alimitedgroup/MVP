package model

type GoodID string

type GoodStock struct {
	ID       string
	Quantity int64
}

type GoodInfo struct {
	ID          string
	Name        string
	Description string
}
