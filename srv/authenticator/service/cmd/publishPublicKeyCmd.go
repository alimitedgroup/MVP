package servicecmd

type PublishPublicKeyCmd struct {
	pemPuk *[]byte
	issuer string
}

func NewPublishPublicKeyCmd(pemPuk *[]byte, issuer string) *PublishPublicKeyCmd {
	return &PublishPublicKeyCmd{pemPuk: pemPuk, issuer: issuer}
}

func (ppkc *PublishPublicKeyCmd) GetKey() *[]byte {
	return ppkc.pemPuk
}

func (ppkc *PublishPublicKeyCmd) GetIssuer() string {
	return ppkc.issuer
}
