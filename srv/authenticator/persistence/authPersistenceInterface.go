package persistence

type IAuthPersistance interface {
	StorePemKeyPair(prk []byte, puk []byte) error
	GetPemPublicKey() (PemPublicKey, error)
	GetPemPrivateKey() (PemPrivateKey, error)
	CheckKeyPairExistence() error
}
