package adapterin

import (
	"testing"

	"github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestGetQuery(t *testing.T) {
	ts := start(t)

	ts.queryResult.EXPECT().GetQueryRule(gomock.Any()).Return(types.QueryRule{
		GoodId:    "1",
		Operator:  "<",
		Threshold: 10,
	}, nil)

	resp, err := ts.nc.Request("notification.queries.get", []byte(`123e4567-e89b-12d3-a456-426614174000`), nats.DefaultTimeout)
	require.NoError(t, err)
	require.JSONEq(t, `{"GoodId":"1","Operator":"<","Threshold":10}`, string(resp.Data))
}
