package serviceportout

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/authenticator/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/authenticator/service/response"
)

type GetPemPrivateKeyPort interface {
	GetPemPrivateKey(cmd *servicecmd.GetPemPrivateKeyCmd) *serviceresponse.GetPemPrivateKeyResponse
}
