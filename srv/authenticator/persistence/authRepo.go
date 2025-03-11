package persistence

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"reflect"

	common "github.com/alimitedgroup/MVP/srv/authenticator/authCommon"
)

type AuthRepository struct {
	prk *PemPrivateKey
	puk *PemPublicKey
}

func NewAuthRepo() *AuthRepository {
	return &AuthRepository{prk: nil, puk: nil}
}

func (ar *AuthRepository) checkKeyPair(prk *[]byte, puk *[]byte) bool {
	decodedPrk, _ := pem.Decode(*prk)
	decodedPuk, _ := pem.Decode(*puk)
	if decodedPrk == nil || decodedPuk == nil {
		return false
	}
	_, errPrk := x509.ParseECPrivateKey(decodedPrk.Bytes)
	pukDecoded, errPuk := x509.ParsePKIXPublicKey(decodedPuk.Bytes)
	if errPrk == nil && errPuk == nil && reflect.TypeOf(pukDecoded) == reflect.TypeOf(&ecdsa.PublicKey{}) {
		return true
	}
	return false
}

func (ar *AuthRepository) StoreKeyPair(prk []byte, puk []byte) error {
	//Store key in PEM format to memory
	if ar.checkKeyPair(&prk, &puk) {
		ar.prk = NewPemPrivateKey(&prk)
		ar.puk = NewPemPublicKey(&puk)
		return nil
	}
	return common.ErrKeyPairNotValid
}

func (ar *AuthRepository) GetPublicKey() (PemPublicKey, error) {
	if ar.puk != nil && len(ar.puk.GetBytes()) > 0 {
		return *ar.puk, nil
	}
	return *NewPemPublicKey(nil), common.ErrNoPublicKey
}

func (ar *AuthRepository) GetPrivateKey() (PemPrivateKey, error) {
	if ar.puk != nil && len(ar.puk.GetBytes()) > 0 {
		return *ar.prk, nil
	}
	return *NewPemPrivateKey(nil), common.ErrNoPublicKey
}

func (ar *AuthRepository) CheckKeyPair() error {
	if ar.prk != nil && ar.puk != nil {
		return nil
	}
	return common.ErrNoKeyPair
}
