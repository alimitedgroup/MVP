package service_Cmd

import "github.com/alimitedgroup/MVP/common/stream"

type MultipleGoodsQuantityCmd struct {
	warehouseID string
	goods       []stream.StockUpdateGood
}

func NewMultipleGoodsQuantityCmd(warehouseID string, goods []stream.StockUpdateGood) *MultipleGoodsQuantityCmd {
	return &MultipleGoodsQuantityCmd{warehouseID: warehouseID, goods: goods}
}

func (mgqc *MultipleGoodsQuantityCmd) GetGoods() []stream.StockUpdateGood {
	return mgqc.goods
}

func (mgqc *MultipleGoodsQuantityCmd) GetWarehouseID() string {
	return mgqc.warehouseID
}
