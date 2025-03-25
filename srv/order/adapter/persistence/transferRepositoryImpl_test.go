package persistence

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferRepositoryImpl(t *testing.T) {
	repo := NewTransferRepositoryImpl()

	_, err := repo.GetTransfer("1")
	require.Error(t, err)
	require.Equal(t, ErrTransferNotFound, err)

	require.False(t, repo.SetTransfer("2", Transfer{
		ID:                "2",
		Status:            "Created",
		LinkedStockUpdate: 0,
	}))
	transfer2, err := repo.GetTransfer("2")
	require.NoError(t, err)
	require.Equal(t, transfer2.ID, "2")

	transfers := repo.GetTransfers()
	require.Len(t, transfers, 1)
	require.Equal(t, transfers[0].ID, "2")

	err = repo.IncrementLinkedStockUpdate("2")
	require.NoError(t, err)

	err = repo.IncrementLinkedStockUpdate("1")
	require.Error(t, err)
	require.Equal(t, ErrTransferNotFound, err)

	require.NoError(t, repo.SetComplete("2"))
	require.Error(t, repo.SetComplete("1"))
}
