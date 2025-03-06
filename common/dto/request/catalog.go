package request

import "github.com/alimitedgroup/MVP/common/stream"

type GetGoodsInfoDTO struct{}

type GetWarehousesInfoDTO struct{}

type AddChangeGoodDTO struct { //jetstream
	Id          string
	Name        string
	Description string
}

type GetGoodsQuantityDTO struct{}

type SetMultipleGoodsQuantityDTO struct { //jetstream
	WarehouseID string
	Goods       []stream.StockUpdateGood
}

type SetGoodQuantityDTO struct { //non utilizzato da controller
	WarehouseId string
	GoodId      string
	NewQuantity int64
}
