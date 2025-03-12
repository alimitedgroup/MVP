package adapter

import (
	"sync"
	"testing"
	"unsafe"

	common "github.com/alimitedgroup/MVP/srv/authenticator/authCommon"
	persistence "github.com/alimitedgroup/MVP/srv/authenticator/persistence"
	servicecmd "github.com/alimitedgroup/MVP/srv/authenticator/service/cmd"
	"github.com/magiconair/properties/assert"
	"go.uber.org/fx"
)

//INIZIO MOCK REPO

var (
	returnErr bool
	mutex     sync.Mutex
)

func setReturnErr(val bool) {
	mutex.Lock()
	returnErr = val
	mutex.Unlock()
}

func getReturnErr() bool {
	mutex.Lock()
	result := returnErr
	mutex.Unlock()
	return result
}

type FakeRepo struct{}

func NewFakeRepo() *FakeRepo {
	return &FakeRepo{}
}

func (fr *FakeRepo) StorePemKeyPair(prk []byte, puk []byte) error {
	if len(prk) == 0 || len(puk) == 0 {
		return common.ErrKeyPairNotValid
	}
	return nil
}

func (fr *FakeRepo) GetPemPublicKey() (persistence.PemPublicKey, error) {
	returnErr := getReturnErr()
	if returnErr {
		return *persistence.NewPemPublicKey(nil), common.ErrNoPublicKey
	}
	returnSt := []byte("test")
	return *persistence.NewPemPublicKey(&returnSt), nil
}

func (fr *FakeRepo) GetPemPrivateKey() (persistence.PemPrivateKey, error) {
	returnErr := getReturnErr()
	if returnErr {
		return *persistence.NewPemPrivateKey(nil), common.ErrNoPrivateKey
	}
	returnSt := []byte("test")
	return *persistence.NewPemPrivateKey(&returnSt), nil
}

func (fr *FakeRepo) CheckKeyPairExistance() error {
	returnErr := getReturnErr()
	if returnErr {
		return common.ErrNoKeyPair
	}
	return nil
}

var p = fx.Options(
	fx.Provide(
		fx.Annotate(NewFakeRepo,
			fx.As(new(persistence.IAuthPersistance)),
		)))

//FINE MOCK REPO

func TestStorePemKeyPair(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthAdapter),
		p,
		fx.Invoke(func(ar *AuthAdapter) {
			fakeKey := []byte("test")
			response := ar.StorePemKeyPair(servicecmd.NewStorePemKeyPairCmd(&fakeKey, &fakeKey))
			assert.Equal(t, response.GetError(), nil)
		}),
	)
}

func TestStoreWrongPemKeyPair(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthAdapter),
		p,
		fx.Invoke(func(ar *AuthAdapter) {
			fakeKey := []byte("test")
			response := ar.StorePemKeyPair(servicecmd.NewStorePemKeyPairCmd(nil, &fakeKey))
			assert.Equal(t, response.GetError(), common.ErrKeyPairNotValid)
		}),
	)
}

func TestGetPemPublicKey(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthAdapter),
		p,
		fx.Invoke(func(ar *AuthAdapter) {
			setReturnErr(false)
			response := ar.GetPemPublicKey(servicecmd.NewGetPemPublicKeyCmd())
			assert.Equal(t, *response.GetPemPublicKey(), []byte("test"))
		}),
	)
}

func TestGetPemPublicKey_NoKey(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthAdapter),
		p,
		fx.Invoke(func(ar *AuthAdapter) {
			setReturnErr(true)
			response := ar.GetPemPublicKey(servicecmd.NewGetPemPublicKeyCmd())
			assert.Equal(t, response.GetPemPublicKey(), (*[]byte)(unsafe.Pointer(nil)))
		}),
	)
}

func TestGetPemPrivateKey(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthAdapter),
		p,
		fx.Invoke(func(ar *AuthAdapter) {
			setReturnErr(false)
			response := ar.GetPemPrivateKey(servicecmd.NewGetPemPrivateKeyCmd())
			assert.Equal(t, *response.GetPemPrivateKey(), []byte("test"))
		}),
	)
}

func TestGetPemPrivateKey_NoKey(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthAdapter),
		p,
		fx.Invoke(func(ar *AuthAdapter) {
			setReturnErr(true)
			response := ar.GetPemPrivateKey(servicecmd.NewGetPemPrivateKeyCmd())
			assert.Equal(t, response.GetPemPrivateKey(), (*[]byte)(unsafe.Pointer(nil)))
		}),
	)
}

func TestCheckKeyPairExistence(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthAdapter),
		p,
		fx.Invoke(func(ar *AuthAdapter) {
			setReturnErr(false)
			response := ar.CheckKeyPairExistance(servicecmd.NewCheckPemKeyPairExistence())
			assert.Equal(t, response.GetError(), nil)
		}),
	)
}

func TestCheckKeyPairExistence_NoKey(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthAdapter),
		p,
		fx.Invoke(func(ar *AuthAdapter) {
			setReturnErr(true)
			response := ar.CheckKeyPairExistance(servicecmd.NewCheckPemKeyPairExistence())
			assert.Equal(t, response.GetError(), common.ErrNoKeyPair)
		}),
	)
}
