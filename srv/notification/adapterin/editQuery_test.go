package adapterin

import (
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestEditQuery(t *testing.T) {
	ts := start(t)

	ts.queryResult.EXPECT().EditQueryRule(gomock.Any(), gomock.Any()).Return(nil)

	resp, err := ts.nc.Request("notification.queries.edit", []byte(`{
		"id": "123e4567-e89b-12d3-a456-426614174000",
		"good_id": "1",
		"operator": "<",
		"threshold": 10
	}`), nats.DefaultTimeout)
	require.NoError(t, err)
	require.Equal(t, "OK", string(resp.Data))
}
