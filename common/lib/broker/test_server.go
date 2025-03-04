package broker

import (
	"os"
	"testing"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

// NewInProcessNATSServer creates a temporary, in-process NATS server, and returns a connection to it. Useful for tests
func NewInProcessNATSServer(t *testing.T) *nats.Conn {
	tmp, err := os.MkdirTemp("", "nats_test")
	require.NoError(t, err)

	server, err := natsserver.NewServer(&natsserver.Options{
		DontListen: true,
		JetStream:  true,
		StoreDir:   tmp,
	})
	require.NoError(t, err)

	server.Start()
	t.Cleanup(func() {
		server.Shutdown()
		err := os.RemoveAll(tmp)
		require.NoError(t, err)
	})
	require.True(t, server.ReadyForConnections(1*time.Second))

	// Create a connection.
	conn, err := nats.Connect("", nats.InProcessServer(server))
	require.NoError(t, err)

	return conn
}
