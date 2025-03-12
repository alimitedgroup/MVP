package persistence

type PemPrivateKey struct {
	prk    *[]byte
	issuer string
}

func NewPemPrivateKey(puk *[]byte, issuer string) *PemPrivateKey {
	return &PemPrivateKey{prk: puk, issuer: issuer}
}

func (prk *PemPrivateKey) GetIssuer() string {
	return prk.issuer
}

func (puk *PemPrivateKey) GetBytes() []byte {
	return *puk.prk
}
