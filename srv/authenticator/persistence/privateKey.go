package persistence

type PemPrivateKey struct {
	prk    *[]byte
	issuer string
}

func NewPemPrivateKey(prk *[]byte, issuer string) *PemPrivateKey {
	return &PemPrivateKey{prk: prk, issuer: issuer}
}

func (prk *PemPrivateKey) GetIssuer() string {
	return prk.issuer
}

func (prk *PemPrivateKey) GetBytes() []byte {
	return *prk.prk
}
