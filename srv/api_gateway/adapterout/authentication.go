package adapterout

import (
	"encoding/json"
	"fmt"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portout"
	"github.com/nats-io/nats.go"
)

type AuthenticationAdapter struct {
	Broker *broker.NatsMessageBroker
}

func (aa *AuthenticationAdapter) GetToken(username string) (types.UserToken, error) {
	req := dto.AuthLoginRequest{Username: username}
	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := aa.Broker.Nats.Request("auth.login", body, nats.DefaultTimeout)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	var respBody dto.AuthLoginResponse
	err = json.Unmarshal(resp.Data, &respBody)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return types.UserToken(respBody.Token), nil
}

func (*AuthenticationAdapter) GetRole(token types.UserToken) (types.UserRole, error) {
	return 0, nil
}

func NewAuthenticationAdapter(broker *broker.NatsMessageBroker) portout.AuthenticationPortOut {
	return &AuthenticationAdapter{Broker: broker}
}
