package model

type GoodID string

type GoodStock struct {
	GoodID   string
	Quantity int64
}

type GoodInfo struct {
	GoodID      string
	Name        string
	Description string
}
