package persistence

type IAuthPersistance interface {
	StorePemKeyPair(prk []byte, puk []byte, emit string) error
	GetPemPublicKey() (PemPublicKey, error)
	GetPemPrivateKey() (PemPrivateKey, error)
	CheckKeyPairExistence() error
}
