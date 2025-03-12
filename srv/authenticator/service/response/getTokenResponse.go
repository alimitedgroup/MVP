package serviceresponse

type GetTokenResponse struct {
	token string
	err   error
}

func NewGetTokenResponse(token string, err error) *GetTokenResponse {
	return &GetTokenResponse{token: token, err: err}
}

func (gtk *GetTokenResponse) GetToken() string {
	return gtk.token
}

func (gtk *GetTokenResponse) GetError() error {
	return gtk.err
}
