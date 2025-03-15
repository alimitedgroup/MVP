package port

import (
	"context"
)

type ICalculateAvailabilityUseCase interface {
	GetAvailable(context.Context, CalculateAvailabilityCmd) (CalculateAvailabilityResponse, error)
}

type CalculateAvailabilityCmd struct {
	Goods []RequestedGood
}

type RequestedGood struct {
	GoodID   string
	Quantity int64
}

type CalculateAvailabilityResponse struct {
	Warehouses []WarehouseAvailability
}

type WarehouseAvailability struct {
	WarehouseID string
	Goods       map[string]int64
}
