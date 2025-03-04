package persistence

type StockRepository interface {
	GetStock(string string) int64
	SetStock(string string, stock int64) bool
	AddStock(string string, stock int64) bool
}
