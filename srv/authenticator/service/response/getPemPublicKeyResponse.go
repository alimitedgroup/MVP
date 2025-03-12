package serviceresponse

type GetPemPublicKeyResponse struct {
	puk *[]byte
	err error
}

func NewGetPemPublicKeyResponse(puk *[]byte, err error) *GetPemPublicKeyResponse {
	return &GetPemPublicKeyResponse{puk: puk, err: err}
}

func (gppkr *GetPemPublicKeyResponse) GetPemPublicKey() *[]byte {
	return gppkr.puk
}
