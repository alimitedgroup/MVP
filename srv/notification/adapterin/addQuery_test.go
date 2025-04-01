package adapterin

import (
	"testing"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestAddQuery(t *testing.T) {
	ts := start(t)

	ts.queryResult.EXPECT().AddQueryRule(gomock.Any()).Return(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), nil)

	resp, err := ts.nc.Request("notification.queries.add", []byte(`{
		"good_id": "1",
		"operator": "<",
		"threshold": 10
	}`), nats.DefaultTimeout)
	require.NoError(t, err)
	require.Equal(t, "123e4567-e89b-12d3-a456-426614174000", string(resp.Data))
}
