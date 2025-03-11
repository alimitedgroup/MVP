package serviceresponse

import "github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"

type GetWarehousesResponse struct {
	warehouseMap map[string]catalogCommon.Warehouse
}

func NewGetWarehousesResponse(warehouseMap map[string]catalogCommon.Warehouse) *GetWarehousesResponse {
	return &GetWarehousesResponse{warehouseMap: warehouseMap}
}

func (gwr *GetWarehousesResponse) GetWarehouseMap() map[string]catalogCommon.Warehouse {
	return gwr.warehouseMap
}
