package serviceportout

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/authenticator/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/authenticator/service/response"
)

type CheckKeyPairExistance interface {
	CheckKeyPairExistance(cmd *servicecmd.CheckPemKeyPairExistence) *serviceresponse.CheckKeyPairExistenceResponse
}
