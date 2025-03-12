package serviceresponse

type GetPemPublicKeyResponse struct {
	puk    *[]byte
	issuer string
	err    error
}

func NewGetPemPublicKeyResponse(puk *[]byte, issuer string, err error) *GetPemPublicKeyResponse {
	return &GetPemPublicKeyResponse{puk: puk, issuer: issuer, err: err}
}

func (gppkr *GetPemPublicKeyResponse) GetIssuer() string {
	return gppkr.issuer
}

func (gppkr *GetPemPublicKeyResponse) GetPemPublicKey() *[]byte {
	return gppkr.puk
}

func (gppkr *GetPemPublicKeyResponse) GetError() error {
	return gppkr.err
}
