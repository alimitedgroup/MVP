package persistence

type StockRepository interface {
	GetStock(goodId string) int64
	SetStock(goodId string, stock int64) bool
	AddStock(goodId string, stock int64) bool
}
