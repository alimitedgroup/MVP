package adapterin

import (
	"testing"

	"github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

func TestListQuery(t *testing.T) {
	ts := start(t)

	ts.queryResult.EXPECT().ListQueryRules().Return([]types.QueryRuleWithId{
		{
			RuleId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			QueryRule: types.QueryRule{
				GoodId:    "1",
				Operator:  "<",
				Threshold: 10,
			},
		},
	}, nil)

	resp, err := ts.nc.Request("notification.queries.list", []byte(""), nats.DefaultTimeout)
	require.NoError(t, err)
	require.JSONEq(t, `[{"RuleId":"123e4567-e89b-12d3-a456-426614174000","GoodId":"1","Operator":"<","Threshold":10}]`, string(resp.Data))
}
