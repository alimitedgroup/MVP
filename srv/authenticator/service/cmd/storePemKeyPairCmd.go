package servicecmd

type StorePemKeyPairCmd struct {
	prk    *[]byte
	puk    *[]byte
	issuer string
}

func NewStorePemKeyPairCmd(prk *[]byte, puk *[]byte, issuer string) *StorePemKeyPairCmd {
	return &StorePemKeyPairCmd{prk: prk, puk: puk, issuer: issuer}
}

func (skpc *StorePemKeyPairCmd) GetPemPrivateKey() *[]byte {
	return skpc.prk
}

func (skpc *StorePemKeyPairCmd) GetPemPublicKey() *[]byte {
	return skpc.puk
}

func (skpc *StorePemKeyPairCmd) GetIssuer() string {
	return skpc.issuer
}
