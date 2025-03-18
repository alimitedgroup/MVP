package persistence

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStockRepositoryImplGetAndSet(t *testing.T) {
	repo := NewStockRepositoryImpl()

	require.False(t, repo.SetStock("1", 0))
	require.True(t, repo.AddStock("1", 10))

	require.Equal(t, repo.GetStock("1"), int64(10))
	require.Equal(t, repo.GetStock("2"), int64(0))

	require.Nil(t, repo.ReserveStock("1", "1", 5))
	require.Nil(t, repo.UnReserveStock("1", 5))
}

func TestStockRepositoryImplReserveAndUnreserve(t *testing.T) {
	repo := NewStockRepositoryImpl()

	require.False(t, repo.SetStock("1", 10))
	require.Nil(t, repo.ReserveStock("1", "1", 5))
	require.Nil(t, repo.UnReserveStock("1", 5))
}

func TestStockRepositoryImplReserveNotExistingGood(t *testing.T) {
	repo := NewStockRepositoryImpl()

	require.NotNil(t, repo.ReserveStock("1", "1", 5))
	require.NotNil(t, repo.UnReserveStock("1", 5))
}
