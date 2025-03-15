package serviceresponse

type StorePemKeyPairResponse struct {
	err error
}

func NewStorePemKeyPairResponse(err error) *StorePemKeyPairResponse {
	return &StorePemKeyPairResponse{err: err}
}

func (spkpr *StorePemKeyPairResponse) GetError() error {
	return spkpr.err
}
