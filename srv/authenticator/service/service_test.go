package service

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"sync"
	"testing"
	"time"

	common "github.com/alimitedgroup/MVP/srv/authenticator/authCommon"
	servicecmd "github.com/alimitedgroup/MVP/srv/authenticator/service/cmd"
	serviceobject "github.com/alimitedgroup/MVP/srv/authenticator/service/object"
	serviceportout "github.com/alimitedgroup/MVP/srv/authenticator/service/portOut"
	serviceresponse "github.com/alimitedgroup/MVP/srv/authenticator/service/response"
	serviceauthenticator "github.com/alimitedgroup/MVP/srv/authenticator/service/strategy"
	"github.com/golang-jwt/jwt/v5"
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

//INIZIO MOCK PORTE

var (
	mutex sync.Mutex
	puk   *[]byte
)

func setPuk(value *[]byte) {
	mutex.Lock()
	defer mutex.Unlock()
	puk = value
}

func getPuk() *[]byte {
	mutex.Lock()
	response := puk
	mutex.Unlock()
	return response
}

type FakeAdapter struct {
	prk *[]byte
	puk *[]byte
}

func NewFakeAdapter() *FakeAdapter {
	return &FakeAdapter{prk: nil, puk: nil}
}

func (fa *FakeAdapter) CheckKeyPairExistance(cmd *servicecmd.CheckPemKeyPairExistenceCmd) *serviceresponse.CheckKeyPairExistenceResponse {
	if fa.prk == nil || fa.puk == nil {
		return serviceresponse.NewCheckKeyPairExistenceResponse(common.ErrNoKeyPair)
	}
	return serviceresponse.NewCheckKeyPairExistenceResponse(nil)
}

func (fa *FakeAdapter) GetPemPrivateKey(cmd *servicecmd.GetPemPrivateKeyCmd) *serviceresponse.GetPemPrivateKeyResponse {
	return serviceresponse.NewGetPemPrivateKeyResponse(fa.prk, "test-issuer", nil)
}

func (fa *FakeAdapter) GetPemPublicKey(cmd *servicecmd.GetPemPublicKeyCmd) *serviceresponse.GetPemPublicKeyResponse {
	return serviceresponse.NewGetPemPublicKeyResponse(fa.puk, "test-issuer", nil)
}

func (fa *FakeAdapter) StorePemKeyPair(cmd *servicecmd.StorePemKeyPairCmd) *serviceresponse.StorePemKeyPairResponse {
	fa.prk = cmd.GetPemPrivateKey()
	fa.puk = cmd.GetPemPublicKey()
	return serviceresponse.NewStorePemKeyPairResponse(nil)
}

func (fa *FakeAdapter) PublishKey(cmd *servicecmd.PublishPublicKeyCmd) *serviceresponse.PublishPublicKeyResponse {
	setPuk(cmd.GetKey())
	return serviceresponse.NewPublishPublicKeyResponse(nil)
}

func (fa *FakeAdapter) Authenticate(us serviceobject.UserData) (string, error) {
	return "test-role", nil
}

var p = fx.Option(
	fx.Provide(
		NewAuthService,
		fx.Annotate(NewFakeAdapter,
			fx.As(new(serviceportout.ICheckKeyPairExistance)),
			fx.As(new(serviceportout.IGetPemPrivateKeyPort)),
			fx.As(new(serviceportout.IGetPemPublicKeyPort)),
			fx.As(new(serviceportout.IStorePemKeyPair)),
			fx.As(new(serviceportout.IPublishPort)),
			fx.As(new(serviceauthenticator.IAuthenticateUserStrategy)),
		),
	),
)

//FINE MOCK PORT

func TestGetToken(t *testing.T) {
	fx.New(
		p,
		fx.Invoke(func(as *AuthService) {
			tokenResponse := as.GetToken(servicecmd.NewGetTokenCmd("test-username"))
			time.Sleep(1 * time.Second)
			pemKey := getPuk()
			decodedPuk, _ := pem.Decode(*pemKey)
			if decodedPuk == nil {
				t.Error("cannot decode Public Key")
				return
			}
			puk, errPuk := x509.ParsePKIXPublicKey(decodedPuk.Bytes)
			if errPuk != nil {
				t.Error("cannot decode Public Key")
			}
			claims := jwt.MapClaims{}
			token, _ := jwt.ParseWithClaims(tokenResponse.GetToken(), &claims, func(token *jwt.Token) (interface{}, error) {
				return puk, nil
			})
			assert.Equal(t, token.Valid, true)
		}),
	)
}
