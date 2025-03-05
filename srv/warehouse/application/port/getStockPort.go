package port

type GetStockPort interface {
	GetStock(goodId string) int64
}
