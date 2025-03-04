package goodRepository

type IGoodRepository interface {
	GetGoods() map[string]good
	GetGoodsGlobalQuantity() map[string]int64
	SetGoodQuantity(warehouseID string, goodID string, newQuantity int64) error
	addWarehouse(warehouseID string)
	AddGood(name string, description string, goodID string) error
	ChangeGoodData(goodID string, newName string, newDescription string) error
}
