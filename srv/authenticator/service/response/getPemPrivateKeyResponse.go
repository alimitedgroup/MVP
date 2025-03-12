package serviceresponse

type GetPemPrivateKeyResponse struct {
	prk *[]byte
	err error
}

func NewGetPemPrivateKeyResponse(prk *[]byte, err error) *GetPemPrivateKeyResponse {
	return &GetPemPrivateKeyResponse{prk: prk, err: err}
}

func (gppkr *GetPemPrivateKeyResponse) GetPemPrivateKey() *[]byte {
	return gppkr.prk
}
