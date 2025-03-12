package serviceresponse

type PublishPublicKeyResponse struct {
	err error
}

func NewPublishPublicKeyResponse(err error) *PublishPublicKeyResponse {
	return &PublishPublicKeyResponse{err: err}
}

func (ppkr *PublishPublicKeyResponse) GetErrror() error {
	return ppkr.err
}
