package serviceportout

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/authenticator/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/authenticator/service/response"
)

type ICheckKeyPairExistance interface {
	CheckKeyPairExistance(cmd *servicecmd.CheckPemKeyPairExistenceCmd) *serviceresponse.CheckKeyPairExistenceResponse
}
