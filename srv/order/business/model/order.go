package model

type OrderID string

type Order struct {
	Id           OrderID
	Status       string
	UpdateTime   int64
	CreationTime int64
	Name         string
	FullName     string
	Address      string
	Goods        []GoodStock
	Reservations []string
	Warehouses   []OrderWarehouseUsed
}

type OrderWarehouseUsed struct {
	WarehouseID string
	Goods       map[GoodID]int64
}

func (o *Order) IsCompleted() bool {
	m := make(map[string]int64)

	for _, good := range o.Goods {
		old, exist := m[string(good.ID)]
		if !exist {
			old = 0
		}
		m[string(good.ID)] = old + good.Quantity
	}

	for _, warehouse := range o.Warehouses {
		for goodId, quantity := range warehouse.Goods {
			m[string(goodId)] -= quantity
		}
	}

	for _, quantity := range m {
		if quantity != 0 {
			return false
		}
	}

	return true
}
