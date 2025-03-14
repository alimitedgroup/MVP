package model

type OrderID string

type Order struct {
	Id           OrderID
	Status       string
	CreationTime int64
	Name         string
	Email        string
	Address      string
	Goods        []GoodStock
}
