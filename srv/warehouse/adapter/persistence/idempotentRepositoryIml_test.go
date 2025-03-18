package persistence

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIdempotenteRepositoryImpl(t *testing.T) {
	repo := NewIdempotentRepositoryImpl()
	repo.SaveEventID("event", "id")

	require.True(t, repo.IsAlreadyProcessed("event", "id"))
	require.False(t, repo.IsAlreadyProcessed("event", "id2"))
	require.False(t, repo.IsAlreadyProcessed("event2", "id"))
}
