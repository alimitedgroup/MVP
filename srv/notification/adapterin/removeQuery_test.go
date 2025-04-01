package adapterin

import (
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestRemoveQuery(t *testing.T) {
	ts := start(t)

	ts.queryResult.EXPECT().RemoveQueryRule(gomock.Any()).Return(nil)

	resp, err := ts.nc.Request("notification.queries.remove", []byte(`123e4567-e89b-12d3-a456-426614174000`), nats.DefaultTimeout)
	require.NoError(t, err)
	require.Equal(t, "OK", string(resp.Data))
}
