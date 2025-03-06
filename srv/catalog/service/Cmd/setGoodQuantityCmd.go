package service_Cmd

type SetGoodQuantityCmd struct {
	warehouseId string
	goodId      string
	newQuantity int64
}

func NewSetGoodQuantityCmd(warehouseId string, goodId string, newQuantity int64) *SetGoodQuantityCmd {
	return &SetGoodQuantityCmd{warehouseId: warehouseId, goodId: goodId, newQuantity: newQuantity}
}

func (sgqc *SetGoodQuantityCmd) GetGoodId() string {
	return sgqc.goodId
}

func (sgqc *SetGoodQuantityCmd) GetWarehouseId() string {
	return sgqc.warehouseId
}

func (sgqc *SetGoodQuantityCmd) GetNewQuantity() int64 {
	return sgqc.newQuantity
}
