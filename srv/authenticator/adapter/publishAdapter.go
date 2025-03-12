package adapter

import (
	"crypto/x509"
	"encoding/pem"

	common "github.com/alimitedgroup/MVP/srv/authenticator/authCommon"
	"github.com/alimitedgroup/MVP/srv/authenticator/publisher"
	servicecmd "github.com/alimitedgroup/MVP/srv/authenticator/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/authenticator/service/response"
)

type AuthPublisherAdapter struct {
	pb publisher.IAuthPublisher
}

func NewAuthPublisherAdapter(pb *publisher.AuthPublisher) *AuthPublisherAdapter {
	return &AuthPublisherAdapter{pb: pb}
}

func (apa *AuthPublisherAdapter) PublishKey(cmd *servicecmd.PublishPublicKeyCmd) *serviceresponse.PublishPublicKeyResponse {
	decodedPuk, _ := pem.Decode(*cmd.GetKey())
	if decodedPuk == nil {
		return serviceresponse.NewPublishPublicKeyResponse(common.ErrPublish)
	}
	pukDecoded, errPuk := x509.ParsePKIXPublicKey(decodedPuk.Bytes)
	if errPuk != nil {
		return serviceresponse.NewPublishPublicKeyResponse(errPuk)
	}
	response := apa.pb.PublishKey(pukDecoded, cmd.GetIssuer())
	return serviceresponse.NewPublishPublicKeyResponse(response)
}
