package publisher

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	common "github.com/alimitedgroup/MVP/srv/authenticator/authCommon"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/magiconair/properties/assert"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
)

var (
	mutex     sync.Mutex
	published bool
)

func setPublish(value bool) {
	mutex.Lock()
	defer mutex.Unlock()
	published = value
}

func getPublish() bool {
	mutex.Lock()
	value := published
	mutex.Unlock()
	return value
}

func GeneratePublicKey() (crypto.PublicKey, error) {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	puk := pk.Public()
	return puk, nil
}

func GenerateWrongPublicKey() (crypto.PublicKey, error) {
	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return pk.Public(), nil
}

func JsDetector(ctx context.Context, msg jetstream.Msg) error {
	var key ecdsa.PublicKey
	err := jwk.ParseRawKey(msg.Data(), &key)
	if err != nil {
		setPublish(false)
	}
	setPublish(true)
	return nil
}

func TestPublishing(t *testing.T) {
	ns, _ := broker.NewInProcessNATSServer(t)
	ctx := context.Background()
	app := fx.New(
		fx.Supply(ns),
		fx.Provide(NewPublisher),
		fx.Provide(broker.NewRestoreStreamControl),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Invoke(func(lc fx.Lifecycle, rsc *broker.RestoreStreamControl, pb *AuthPublisher) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					setPublish(false)
					puk, err := GeneratePublicKey()
					if err != nil {
						return err
					}
					err = pb.PublishKey(puk, "test-issuer")
					assert.Equal(t, err, nil)
					err = pb.mb.RegisterJsHandler(ctx, rsc, stream.KeyStream, JsDetector)
					if err != nil {
						return err
					}
					time.Sleep(1 * time.Second)
					assert.Equal(t, getPublish(), true)
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return nil
				},
			})
		}),
	)

	err := app.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = app.Stop(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func TestPublishingWrongKey(t *testing.T) {
	ns, _ := broker.NewInProcessNATSServer(t)
	ctx := context.Background()
	app := fx.New(
		fx.Supply(ns),
		fx.Provide(NewPublisher),
		fx.Provide(broker.NewRestoreStreamControl),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Invoke(func(lc fx.Lifecycle, rsc *broker.RestoreStreamControl, pb *AuthPublisher) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					setPublish(false)
					puk, err := GenerateWrongPublicKey()
					if err != nil {
						return err
					}
					err = pb.PublishKey(puk, "test-issuer")
					assert.Equal(t, err, common.ErrPublish)
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return nil
				},
			})
		}),
	)

	err := app.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = app.Stop(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()
}
