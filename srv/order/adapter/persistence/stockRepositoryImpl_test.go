package persistence

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStockRepositoryImpl(t *testing.T) {
	repo := NewStockRepositoryImpl()

	_, err := repo.GetStock("1", "1")
	require.Error(t, err)
	require.Equal(t, ErrWarehouseNotFound, err)
	require.Equal(t, repo.GetGlobalStock("1"), 0)

	require.False(t, repo.SetStock("1", "2", 1))
	stock2, err := repo.GetStock("1", "2")
	require.NoError(t, err)
	require.Equal(t, stock2, 2)
	require.Equal(t, repo.GetGlobalStock("2"), 1)

	exist, err := repo.AddStock("1", "2", 1)
	require.NoError(t, err)
	require.True(t, exist)

	require.Equal(t, repo.GetWarehouses(), []string{"1"})

}
