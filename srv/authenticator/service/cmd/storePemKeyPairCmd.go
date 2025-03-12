package servicecmd

type StorePemKeyPairCmd struct {
	prk *[]byte
	puk *[]byte
}

func NewStorePemKeyPairCmd(prk *[]byte, puk *[]byte) *StorePemKeyPairCmd {
	return &StorePemKeyPairCmd{prk: prk, puk: puk}
}

func (skpc *StorePemKeyPairCmd) GetPemPrivateKey() *[]byte {
	return skpc.prk
}

func (skpc *StorePemKeyPairCmd) GetPemPublicKey() *[]byte {
	return skpc.puk
}
