package persistence

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"reflect"
	"sync"

	common "github.com/alimitedgroup/MVP/srv/authenticator/authCommon"
)

type AuthRepository struct {
	prk    *PemPrivateKey
	puk    *PemPublicKey
	issuer string
	mutex  sync.Mutex
}

func NewAuthRepo() *AuthRepository {
	return &AuthRepository{prk: nil, puk: nil, issuer: ""}
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

func (ar *AuthRepository) StorePemKeyPair(prk []byte, puk []byte, emit string) error {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	//Store key in PEM format to memory
	if ar.checkKeyPair(&prk, &puk) {
		ar.issuer = emit
		ar.prk = NewPemPrivateKey(&prk, ar.issuer)
		ar.puk = NewPemPublicKey(&puk, ar.issuer)
		return nil
	}
	return common.ErrKeyPairNotValid
}

func (ar *AuthRepository) GetPemPublicKey() (PemPublicKey, error) {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	if ar.puk != nil && len(ar.puk.GetBytes()) > 0 {
		return *ar.puk, nil
	}
	return *NewPemPublicKey(nil, ""), common.ErrNoPublicKey
}

func (ar *AuthRepository) GetPemPrivateKey() (PemPrivateKey, error) {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	if ar.puk != nil && len(ar.puk.GetBytes()) > 0 {
		return *ar.prk, nil
	}
	return *NewPemPrivateKey(nil, ""), common.ErrNoPrivateKey
}

func (ar *AuthRepository) CheckKeyPairExistence() error {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	if ar.prk != nil && ar.puk != nil {
		return nil
	}
	return common.ErrNoKeyPair
}
