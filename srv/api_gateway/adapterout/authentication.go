package adapterout

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portout"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"log/slog"
	"time"
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

func (*AuthenticationAdapter) GetUsername(token types.ParsedToken) (string, error) {
	token2, ok := token.(*jwt.Token)
	if !ok {
		return "", portout.ErrTokenInvalid
	}

	sub, err := token2.Claims.GetSubject()
	if err != nil {
		return "", portout.ErrTokenInvalid
	}

	return sub, nil
}

func (*AuthenticationAdapter) GetRole(token types.ParsedToken) (types.UserRole, error) {
	token2, ok := token.(*jwt.Token)
	if !ok {
		return types.RoleNone, portout.ErrTokenInvalid
	}

	roleraw, ok := token2.Claims.(jwt.MapClaims)["role"]
	if !ok {
		return types.RoleNone, portout.ErrTokenInvalid
	}

	rolestr, ok := roleraw.(string)
	if !ok {
		return types.RoleNone, portout.ErrTokenInvalid
	}

	role := types.RoleFromString(rolestr)
	if role == types.RoleNone {
		return types.RoleNone, portout.ErrTokenInvalid
	}
	return role, nil
}

func (aa *AuthenticationAdapter) VerifyToken(token types.UserToken) (types.ParsedToken, error) {
	parsed, err := jwt.Parse(string(token), func(token *jwt.Token) (interface{}, error) {
		iss, err := token.Claims.GetIssuer()
		if err != nil {
			return nil, portout.ErrTokenInvalid
		}
		if iss == "" {
			return nil, portout.ErrTokenInvalid
		}

		key, err := getValidationKey(context.TODO(), aa, iss)
		if err != nil {
			slog.Error("Error getting JWT validation key", "error", err)
			return nil, portout.ErrTokenInvalid
		}

		return key, nil
	}, jwt.WithValidMethods([]string{"ES256"}))
	if err != nil {
		if !errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			slog.Error("Error parsing token", "error", err)
		}

		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, portout.ErrTokenExpired
		}

		return nil, portout.ErrTokenInvalid
	}

	return parsed, nil
}

// getValidationKey returns a public key that can be used to verify JWTs signed by the given issuer
func getValidationKey(ctx context.Context, aa *AuthenticationAdapter, issuer string) (*ecdsa.PublicKey, error) {
	stream, err := aa.Broker.Js.CreateStream(
		ctx,
		jetstream.StreamConfig{Name: "auth_keys", Subjects: []string{"keys.>"}},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}

	consumer, err := stream.CreateConsumer(ctx, jetstream.ConsumerConfig{
		DeliverPolicy: jetstream.DeliverLastPerSubjectPolicy,
		FilterSubject: fmt.Sprintf("keys.%s", issuer),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	msg, err := consumer.Next(jetstream.FetchMaxWait(time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to receive message: %w", err)
	}

	var key ecdsa.PublicKey
	err = jwk.ParseRawKey(msg.Data(), &key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse key: %w", err)
	}

	return &key, nil
}

func NewAuthenticationAdapter(broker *broker.NatsMessageBroker) portout.AuthenticationPortOut {
	return &AuthenticationAdapter{Broker: broker}
}
