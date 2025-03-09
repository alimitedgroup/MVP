package adapterout

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"log"
	"time"
)

type Users map[string]types.UserRole

var DefaultUsers = Users{
	"admin":       types.RoleGlobalAdmin,
	"local_admin": types.RoleLocalAdmin,
	"client":      types.RoleClient,
}

type AuthMock struct {
	issuer  string
	keypair *ecdsa.PrivateKey
	broker  *broker.NatsMessageBroker
	users   Users
}

func NewAuthMock(issuer string, broker *broker.NatsMessageBroker, users Users) *AuthMock {
	return &AuthMock{
		issuer:  issuer,
		broker:  broker,
		users:   users,
		keypair: genKeypair(broker.Js, issuer),
	}
}

func StartAuthMock(mock *AuthMock, lc fx.Lifecycle) {
	subscription, err := mock.broker.Nats.Subscribe("auth.login", func(msg *nats.Msg) {
		var req dto.AuthLoginRequest
		_ = json.Unmarshal(msg.Data, &req)

		log.Println("received request:", string(msg.Data), req)

		role, ok := mock.users[req.Username]
		if !ok {
			jsonBoh, _ := json.Marshal(response.ResponseDTO[string]{Error: "invalid_credentials", Message: "Invalid credentials"})
			err := msg.Respond(jsonBoh)
			if err != nil {
				panic(err)
			}
			return
		}

		token, err := mock.genJwt(req.Username, role.String())
		if err != nil {
			panic(err)
		}
		jsonBoh, err := json.Marshal(dto.AuthLoginResponse{Token: token})
		if err != nil {
			panic(err)
		}
		err = msg.Respond(jsonBoh)
		if err != nil {
			panic(err)
		}
	})

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return err
		},
		OnStop: func(context.Context) error {
			return subscription.Unsubscribe()
		},
	})
}

func (a *AuthMock) genJwt(username string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub":  username,
		"role": role,
		// expiration in 1 week from now
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iss": a.issuer,
	})
	return token.SignedString(a.keypair)
}

func genKeypair(js jetstream.JetStream, issuer string) *ecdsa.PrivateKey {
	keypair, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	// Serializzazione in formato JWK
	key, err := jwk.New(keypair.PublicKey)
	if err != nil {
		panic(err)
	}
	keybytes, err := json.Marshal(key)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	_, err = js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "auth_keys",
		Subjects: []string{"keys.>"},
	})
	if err != nil {
		panic(err)
	}

	_, err = js.Publish(ctx, fmt.Sprintf("keys.%s", issuer), keybytes)
	if err != nil {
		panic(err)
	}

	return keypair
}
