package portin

import "github.com/alimitedgroup/MVP/common/dto"

type WarehouseOverview struct {
	ID string
}

type Warehouses interface {
	GetWarehouseByID(id int64) (dto.Warehouse, error)
	GetWarehouses() ([]WarehouseOverview, error)
	GetGoods() ([]dto.GoodAndAmount, error)
}
