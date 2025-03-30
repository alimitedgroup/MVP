package adapterout

import (
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portout"
	"github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/nats-io/nats.go"
)

type NotificationsAdapterOut struct {
	Broker *broker.NatsMessageBroker
}

func NewNotificationsAdapter(broker *broker.NatsMessageBroker) portout.NotificationPortOut {
	return &NotificationsAdapterOut{
		Broker: broker,
	}
}

func (c NotificationsAdapterOut) CreateQuery(dto dto.Rule) (string, error) {
	payload, err := json.Marshal(dto)
	if err != nil {
		return "", err
	}

	resp, err := c.Broker.Nats.Request("notification.queries.add", payload, nats.DefaultTimeout)
	if err != nil {
		return "", err
	}

	queryId := string(resp.Data)

	return queryId, err
}

func (c NotificationsAdapterOut) GetQueries() ([]types.QueryRuleWithId, error) {
	resp, err := c.Broker.Nats.Request("notification.queries.list", []byte("{}"), nats.DefaultTimeout)
	if err != nil {
		return nil, err
	}

	var respDto []types.QueryRuleWithId
	err = json.Unmarshal(resp.Data, &respDto)
	if err != nil {
		return nil, err
	}

	return respDto, err
}

var _ portout.NotificationPortOut = (*NotificationsAdapterOut)(nil)
