package serviceportout

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/authenticator/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/authenticator/service/response"
)

type GetPemPublicKeyPort interface {
	GetPemPublicKey(cmd *servicecmd.GetPemPublicKeyCmd) *serviceresponse.GetPemPublicKeyResponse
}
