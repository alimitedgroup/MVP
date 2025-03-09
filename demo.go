package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"time"
)

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("token is invalid")
)

var keypair *ecdsa.PrivateKey
var issuer = "ciao"

func genJwt(username string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub":  username,
		"role": role,
		// expiration: 1 week
		"exp": time.Now().Add(-time.Hour * 24 * 7).Unix(),
	})
	return token.SignedString(keypair)
}

// verifyJwt verifica che il JWT fornito sia valido, ritornando un errore se non lo è.
// Se il jwt è scaduto, ritorna ErrTokenExpired.
// Per qualsiasi altra invalidità, ritorna ErrTokenInvalid.
// Se il token è valido, ritorna nil.
func verifyJwt(token string) error {
	_, err := jwt.Parse(
		token,
		func(token *jwt.Token) (interface{}, error) { return &keypair.PublicKey, nil },
		jwt.WithValidMethods([]string{"ES256"}),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return ErrTokenExpired
		}
		return ErrTokenInvalid
	}

	return nil
}

func genKeypair(js jetstream.JetStream) {
	var err error
	keypair, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
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

}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	js, err := jetstream.New(nc)

	genKeypair(js)

	_, err = nc.Subscribe("auth.login", func(msg *nats.Msg) {
		var req dto.AuthLoginRequest
		_ = json.Unmarshal(msg.Data, &req)
		if req.Username == "admin" {
			token, err := genJwt("admin", "local_admin")
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
		} else {
			jsonBoh, _ := json.Marshal(response.ResponseDTO[string]{Error: "invalid_credentials", Message: "Invalid credentials"})
			err = msg.Respond(jsonBoh)
			if err != nil {
				panic(err)
			}
		}
	})
	if err != nil {
		panic(err)
	}

	select {}
}
