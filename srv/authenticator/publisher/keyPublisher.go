package publisher

import (
	"context"
	"crypto"
	"encoding/json"
	"fmt"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/lestrrat-go/jwx/jwk"
)

type AuthPublisher struct {
	mb *broker.NatsMessageBroker
}

func NewPublisher(mb *broker.NatsMessageBroker) *AuthPublisher {
	return &AuthPublisher{mb: mb}
}

func (ap *AuthPublisher) PublishKey(puk crypto.PublicKey, issuer string) error {
	key, err := jwk.New(puk)
	if err != nil {
		return err
	}
	msg, err := json.Marshal(key)
	if err != nil {
		return err
	}
	ctx := context.Background()
	_, _ = ap.mb.Js.CreateStream(ctx, stream.KeyStream)
	_, err = ap.mb.Js.Publish(ctx, fmt.Sprintf("keys.%s", issuer), msg)
	if err != nil {
		return err
	}
	return nil
}
