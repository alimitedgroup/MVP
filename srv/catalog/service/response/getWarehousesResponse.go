package serviceresponse

import (
	"github.com/alimitedgroup/MVP/common/dto"
)

type GetWarehousesResponse struct {
	warehouseMap map[string]dto.Warehouse
}

func NewGetWarehousesResponse(warehouseMap map[string]dto.Warehouse) *GetWarehousesResponse {
	return &GetWarehousesResponse{warehouseMap: warehouseMap}
}

func (gwr *GetWarehousesResponse) GetWarehouseMap() map[string]dto.Warehouse {
	return gwr.warehouseMap
}
