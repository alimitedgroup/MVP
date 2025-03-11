package persistence

type PemPrivateKey struct {
	prk *[]byte
}

func NewPemPrivateKey(puk *[]byte) *PemPrivateKey {
	return &PemPrivateKey{prk: puk}
}

func (puk *PemPrivateKey) GetBytes() []byte {
	return *puk.prk
}
