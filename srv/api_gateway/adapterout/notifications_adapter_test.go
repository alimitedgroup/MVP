package adapterout

import (
	"testing"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

func TestGetQueries(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)

	sub, err := nc.Subscribe("notification.queries.list", func(msg *nats.Msg) {
		err := msg.Respond([]byte(`[{"RuleId":"b8f6208f-0828-469f-811d-748ffbfd24b6","GoodId":"1","Operator":">","Threshold":10}]`))
		require.NoError(t, err)
	})
	require.NoError(t, err)
	defer func() {
		err := sub.Unsubscribe()
		require.NoError(t, err)
	}()

	brk := broker.NewTest(t, nc)
	order := NewNotificationsAdapter(brk)

	info, err := order.GetQueries()
	require.NoError(t, err)
	require.Len(t, info, 1)
	require.Equal(t, "b8f6208f-0828-469f-811d-748ffbfd24b6", info[0].RuleId.String())
}

func TestCreateQuery(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)

	sub, err := nc.Subscribe("notification.queries.add", func(msg *nats.Msg) {
		err := msg.Respond([]byte(`1`))
		require.NoError(t, err)
	})
	require.NoError(t, err)
	defer func() {
		err := sub.Unsubscribe()
		require.NoError(t, err)
	}()

	brk := broker.NewTest(t, nc)
	order := NewNotificationsAdapter(brk)

	dto := dto.Rule{
		GoodId:    "1",
		Operator:  ">",
		Threshold: 10,
	}
	queryId, err := order.CreateQuery(dto)
	require.NoError(t, err)
	require.Equal(t, "1", queryId)
}
