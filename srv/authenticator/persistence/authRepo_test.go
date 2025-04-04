package persistence

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	common "github.com/alimitedgroup/MVP/srv/authenticator/authCommon"
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

func TestStorePemKeyPair(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthRepo),
		fx.Invoke(func(ar *AuthRepository) {
			prk, puk, err := GeneratePemKey()
			assert.Equal(t, err, nil)
			err2 := ar.StorePemKeyPair(*prk, *puk, "test-issuer")
			assert.Equal(t, err2, nil)
		}),
	)
}

func TestStoreWrongKeyPair(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthRepo),
		fx.Invoke(func(ar *AuthRepository) {
			prk, puk, err := GenerateWrongPemKey()
			assert.Equal(t, err, nil)
			err2 := ar.StorePemKeyPair(*prk, *puk, "test-issuer")
			assert.Equal(t, err2, common.ErrKeyPairNotValid)
		}),
	)
}

func TestStoreGarbageKeyPair(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthRepo),
		fx.Invoke(func(ar *AuthRepository) {
			prk := []byte{}
			prk = append(prk, byte(7))
			puk := []byte{}
			puk = append(puk, byte(7))
			err := ar.StorePemKeyPair(prk, puk, "test-issuer")
			assert.Equal(t, err, common.ErrKeyPairNotValid)
		}),
	)
}

func TestGetPemPublicKey(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthRepo),
		fx.Invoke(func(ar *AuthRepository) {
			prk, puk, err := GeneratePemKey()
			assert.Equal(t, err, nil)
			err2 := ar.StorePemKeyPair(*prk, *puk, "test-issuer")
			assert.Equal(t, err2, nil)
			pukc, err3 := ar.GetPemPublicKey()
			assert.Equal(t, err3, nil)
			assert.Equal(t, pukc.GetBytes(), *puk)
			if pukc.GetIssuer() == "" {
				t.Error("Issuer is <nil>")
			}
		}),
	)
}

func TestGetPemPrivateKey(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthRepo),
		fx.Invoke(func(ar *AuthRepository) {
			prk, puk, err := GeneratePemKey()
			assert.Equal(t, err, nil)
			err2 := ar.StorePemKeyPair(*prk, *puk, "test-issuer")
			assert.Equal(t, err2, nil)
			pukc, err3 := ar.GetPemPrivateKey()
			assert.Equal(t, err3, nil)
			assert.Equal(t, pukc.GetBytes(), *prk)
			if pukc.GetIssuer() == "" {
				t.Error("Issuer is <nil>")
			}
		}),
	)
}

func TestGetPemPublicKeyWithNoKey(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthRepo),
		fx.Invoke(func(ar *AuthRepository) {
			pukc, err3 := ar.GetPemPublicKey()
			assert.Equal(t, err3, common.ErrNoPublicKey)
			assert.Equal(t, &pukc, NewPemPublicKey(nil, ""))
		}),
	)
}

func TestGetPemPrivateKeyWithNoKey(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthRepo),
		fx.Invoke(func(ar *AuthRepository) {
			prkc, err3 := ar.GetPemPrivateKey()
			assert.Equal(t, err3, common.ErrNoPrivateKey)
			assert.Equal(t, &prkc, NewPemPrivateKey(nil, ""))
		}),
	)
}

func TestCheckKeyPair(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthRepo),
		fx.Invoke(func(ar *AuthRepository) {
			prk, puk, err := GeneratePemKey()
			assert.Equal(t, err, nil)
			err2 := ar.StorePemKeyPair(*prk, *puk, "test-issuer")
			assert.Equal(t, err2, nil)
			assert.Equal(t, ar.CheckKeyPairExistence(), nil)
		}),
	)
}

func TestCheckKeyPairWithNoKey(t *testing.T) {
	fx.New(
		fx.Provide(NewAuthRepo),
		fx.Invoke(func(ar *AuthRepository) {
			assert.Equal(t, ar.CheckKeyPairExistence(), common.ErrNoKeyPair)
		}),
	)
}
