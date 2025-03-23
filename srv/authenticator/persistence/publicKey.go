package persistence

type PemPublicKey struct {
	puk    *[]byte
	issuer string
}

func NewPemPublicKey(puk *[]byte, issuer string) *PemPublicKey {
	return &PemPublicKey{puk: puk, issuer: issuer}
}

func (puk *PemPublicKey) GetIssuer() string {
	return puk.issuer
}

func (puk *PemPublicKey) GetBytes() []byte {
	return *puk.puk
}
