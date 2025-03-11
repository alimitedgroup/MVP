package persistence

type PemPublicKey struct {
	puk *[]byte
}

func NewPemPublicKey(puk *[]byte) *PemPublicKey {
	return &PemPublicKey{puk: puk}
}

func (puk *PemPublicKey) GetBytes() []byte {
	return *puk.puk
}
