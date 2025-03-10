package response

import "github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"

type GetGoodsDataResponseDTO struct {
	GoodMap map[string]catalogCommon.Good `json:"goodMap"`
	Err     string                        `json:"Error"`
}

type GetWarehouseResponseDTO struct {
	WarehouseMap map[string]catalogCommon.Warehouse `json:"warehouseMap"`
	Err          string                             `json:"Error"`
}

type GetGoodsQuantityResponseDTO struct {
	GoodMap map[string]int64 `json:"goodMap"`
	Err     string           `json:"Error"`
}
