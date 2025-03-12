package serviceresponse

type GetPemPrivateKeyResponse struct {
	prk    *[]byte
	issuer string
	err    error
}

func NewGetPemPrivateKeyResponse(prk *[]byte, issuer string, err error) *GetPemPrivateKeyResponse {
	return &GetPemPrivateKeyResponse{prk: prk, issuer: issuer, err: err}
}

func (gppkr *GetPemPrivateKeyResponse) GetIssuer() string {
	return gppkr.issuer
}

func (gppkr *GetPemPrivateKeyResponse) GetPemPrivateKey() *[]byte {
	return gppkr.prk
}

func (gppkr *GetPemPrivateKeyResponse) GetError() error {
	return gppkr.err
}
