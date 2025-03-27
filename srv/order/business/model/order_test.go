package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOrderIsComplete(t *testing.T) {
	order := Order{
		ID:           "1",
		Status:       "Filled",
		UpdateTime:   time.Now().UnixMilli(),
		CreationTime: time.Now().UnixMilli(),
		Name:         "Test Order",
		FullName:     "Test User",
		Address:      "123 Main St",
		Goods: []GoodStock{
			{
				GoodID:   "1",
				Quantity: 5,
			},
		},
		Reservations: []string{"1"},
		Warehouses: []OrderWarehouseUsed{
			{
				WarehouseID: "1",
				Goods: map[string]int64{
					"1": 5,
				},
			},
		},
	}

	isCompleted := order.IsCompleted()

	require.True(t, isCompleted)
}
