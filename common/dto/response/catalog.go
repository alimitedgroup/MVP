package response

import "github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"

type GetGoodsDataResponseDTO struct {
	GoodMap map[string]catalogCommon.Good `json:"goods"`
	Err     string                        `json:"error"`
}

type GetWarehouseResponseDTO struct {
	WarehouseMap map[string]catalogCommon.Warehouse `json:"warehouse_map"`
	Err          string                             `json:"error"`
}

type GetGoodsQuantityResponseDTO struct {
	GoodMap map[string]int64 `json:"goods"`
	Err     string           `json:"error"`
}
