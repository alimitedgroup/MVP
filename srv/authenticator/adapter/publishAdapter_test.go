package adapter

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"testing"

	common "github.com/alimitedgroup/MVP/srv/authenticator/authCommon"
	"github.com/alimitedgroup/MVP/srv/authenticator/publisher"
	servicecmd "github.com/alimitedgroup/MVP/srv/authenticator/service/cmd"
	"github.com/magiconair/properties/assert"
	"go.uber.org/fx"
)

func GeneratePemKey() (*[]byte, *[]byte, error) {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	puk := pk.Public()
	privKeyBytes, err := x509.MarshalECPrivateKey(pk)
	if err != nil {
		return nil, nil, err
	}
	privateKeyFile := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privKeyBytes})
	pubkeybytes, err := x509.MarshalPKIXPublicKey(puk)
	if err != nil {
		return nil, nil, err
	}
	publicKeyFile := pem.EncodeToMemory(&pem.Block{Type: "EC PUBLIC KEY", Bytes: pubkeybytes})
	return &privateKeyFile, &publicKeyFile, nil
}

func GenerateWrongPemKey() (*[]byte, *[]byte, error) {
	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	puk := pk.Public()
	privKeyBytes := x509.MarshalPKCS1PrivateKey(pk)
	privateKeyFile := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privKeyBytes})
	pubkeybytes := x509.MarshalPKCS1PublicKey(puk.(*rsa.PublicKey))
	publicKeyFile := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubkeybytes})
	return &privateKeyFile, &publicKeyFile, nil
}

//INIZIO MOCK PUBLISHER

type FakePublisher struct {
}

func NewFakePublisher() *FakePublisher {
	return &FakePublisher{}
}

func (fp *FakePublisher) PublishKey(puk crypto.PublicKey, issuer string) error {
	if issuer == "wrong-issuer" {
		return errors.New("test error")
	}
	return nil
}

//FINE MOCK PUBLISHER

func TestAdaptKey(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthPublisherAdapter,
			fx.Annotate(NewFakePublisher,
				fx.As(new(publisher.IAuthPublisher))),
		),
		fx.Invoke(func(apa *AuthPublisherAdapter) {
			_, puk, err := GeneratePemKey()
			assert.Equal(t, err, nil)
			response := apa.PublishKey(servicecmd.NewPublishPublicKeyCmd(puk, "test-issuer"))
			assert.Equal(t, response.GetError(), nil)
		}),
	)
}

func TestAdaptKeyWithWrongIssuer(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthPublisherAdapter,
			fx.Annotate(NewFakePublisher,
				fx.As(new(publisher.IAuthPublisher))),
		),
		fx.Invoke(func(apa *AuthPublisherAdapter) {
			_, puk, err := GeneratePemKey()
			assert.Equal(t, err, nil)
			response := apa.PublishKey(servicecmd.NewPublishPublicKeyCmd(puk, "wrong-issuer"))
			assert.Equal(t, response.GetError(), errors.New("test error"))
		}),
	)
}

func TestAdaptWithInvalidKey(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthPublisherAdapter,
			fx.Annotate(NewFakePublisher,
				fx.As(new(publisher.IAuthPublisher))),
		),
		fx.Invoke(func(apa *AuthPublisherAdapter) {
			test := []byte{7}
			response := apa.PublishKey(servicecmd.NewPublishPublicKeyCmd(&test, "wrong-issuer"))
			assert.Equal(t, response.GetError(), common.ErrPublish)
		}),
	)
}

func TestAdaptWithWrongKey(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthPublisherAdapter,
			fx.Annotate(NewFakePublisher,
				fx.As(new(publisher.IAuthPublisher))),
		),
		fx.Invoke(func(apa *AuthPublisherAdapter) {
			_, puk, err := GenerateWrongPemKey()
			assert.Equal(t, err, nil)
			response := apa.PublishKey(servicecmd.NewPublishPublicKeyCmd(puk, "wrong-issuer"))
			assert.Equal(t, response.GetError().Error(), "x509: failed to parse public key (use ParsePKCS1PublicKey instead for this key format)")
		}),
	)
}
