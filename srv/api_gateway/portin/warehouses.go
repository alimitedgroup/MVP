package portin

import (
	"context"

	"github.com/alimitedgroup/MVP/common/dto"
)

type WarehouseOverview struct {
	ID string
}

type Warehouses interface {
	GetWarehouseByID(id int64) (dto.Warehouse, error)
	GetWarehouses() ([]WarehouseOverview, error)
	GetGoods() ([]dto.GoodAndAmount, error)
	CreateGood(ctx context.Context, name string, description string) (string, error)
	UpdateGood(ctx context.Context, goodId string, name string, description string) error
	AddStock(warehouseId string, goodId string, quantity int64) error
	RemoveStock(warehouseId string, goodId string, quantity int64) error
}
