package port

type IGetStockPort interface {
	GetStock(goodId string) int64
}
