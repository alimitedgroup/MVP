package business

import (
	"context"
	"math/rand"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
)

type SimpleCalculateAvailabilityService struct {
	getStockPort port.IGetStockPort
}

func NewSimpleCalculateAvailabilityService(getStockPort port.IGetStockPort) *SimpleCalculateAvailabilityService {
	return &SimpleCalculateAvailabilityService{getStockPort}
}

func (s *SimpleCalculateAvailabilityService) GetAvailable(ctx context.Context, cmd port.CalculateAvailabilityCmd) (port.CalculateAvailabilityResponse, error) {
	// quick return if global stock counter isn't enough
	total := int64(0)
	reqGoods := make(map[string]int64)
	for _, good := range cmd.Goods {
		if s.getStockPort.GetGlobalStock(model.GoodID(good.GoodID)).Quantity < good.Quantity {
			return port.CalculateAvailabilityResponse{}, port.ErrNotEnoughStock
		}
		reqGoods[good.GoodID] = good.Quantity
		total += good.Quantity
	}

	// shuffle warehouses
	warehouses := s.getStockPort.GetWarehouses()
	for i := range warehouses {
		j := rand.Intn(i + 1)
		warehouses[i], warehouses[j] = warehouses[j], warehouses[i]
	}

	availabilities := make([]port.WarehouseAvailability, 0, len(warehouses))

	// NOTE: disable excluded warehoues for now

	// excluded := make(map[string]struct{})
	// for _, warehouseID := range cmd.ExcludedWarehouses {
	// 	excluded[warehouseID] = struct{}{}
	// }

	for _, warehouse := range warehouses {
		// if _, ok := excluded[warehouse.ID]; ok {
		// 	continue
		// }

		warehouseTotal := int64(0)
		toReserveGoods := make(map[string]int64)

		for goodID, quantity := range reqGoods {
			stock, err := s.getStockPort.GetStock(port.GetStockCmd{
				WarehouseID: warehouse.ID,
				GoodID:      goodID,
			})
			if err != nil {
				continue
			}

			toReserveQty := min(stock.Quantity, quantity)
			toReserveGoods[goodID] = toReserveQty
			reqGoods[goodID] -= toReserveQty
			total -= toReserveQty
			warehouseTotal += toReserveQty
		}

		if warehouseTotal > 0 {
			availabilities = append(availabilities, port.WarehouseAvailability{
				WarehouseID: string(warehouse.ID),
				Goods:       toReserveGoods,
			})
		}
	}

	if total > 0 {
		return port.CalculateAvailabilityResponse{}, port.ErrNotEnoughStock
	}

	return port.CalculateAvailabilityResponse{
		Warehouses: availabilities,
	}, nil
}
