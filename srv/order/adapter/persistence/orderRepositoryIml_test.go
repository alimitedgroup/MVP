package persistence

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrderRepositoryImpl(t *testing.T) {
	repo := NewOrderRepositoryImpl()

	_, err := repo.GetOrder("1")
	require.Error(t, err)
	require.Equal(t, ErrOrderNotFound, err)

	require.False(t, repo.SetOrder("2", Order{
		ID:         "2",
		Status:     "Created",
		Warehouses: []OrderWarehouseUsed{},
	}))
	order2, err := repo.GetOrder("2")
	require.NoError(t, err)
	require.Equal(t, order2.ID, "2")

	orders := repo.GetOrders()
	require.Len(t, orders, 1)
	require.Equal(t, orders[0].ID, "2")

	order2, err = repo.AddCompletedWarehouse("2", "1", map[string]int64{"1": 1})
	require.NoError(t, err)
	require.Equal(t, order2.ID, "2")

	_, err = repo.AddCompletedWarehouse("1", "1", map[string]int64{"1": 2})
	require.Error(t, err)
	require.Equal(t, ErrOrderNotFound, err)

	require.NoError(t, repo.SetComplete("2"))
	require.Error(t, repo.SetComplete("1"))
}
