package serviceresponse

type CheckKeyPairExistenceResponse struct {
	err error
}

func NewCheckKeyPairExistenceResponse(err error) *CheckKeyPairExistenceResponse {
	return &CheckKeyPairExistenceResponse{err: err}
}

func (ckper *CheckKeyPairExistenceResponse) GetError() error {
	return ckper.err
}
