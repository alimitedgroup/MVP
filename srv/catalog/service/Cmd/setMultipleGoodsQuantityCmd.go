package service_Cmd

import "github.com/alimitedgroup/MVP/common/stream"

type SetMultipleGoodsQuantityCmd struct {
	warehouseID string
	goods       []stream.StockUpdateGood
}

func NewSetMultipleGoodsQuantityCmd(warehouseID string, goods []stream.StockUpdateGood) *SetMultipleGoodsQuantityCmd {
	return &SetMultipleGoodsQuantityCmd{warehouseID: warehouseID, goods: goods}
}

func (mgqc *SetMultipleGoodsQuantityCmd) GetGoods() []stream.StockUpdateGood {
	return mgqc.goods
}

func (mgqc *SetMultipleGoodsQuantityCmd) GetWarehouseID() string {
	return mgqc.warehouseID
}
