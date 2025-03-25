package serviceportout

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/authenticator/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/authenticator/service/response"
)

type IStorePemKeyPair interface {
	StorePemKeyPair(cmd *servicecmd.StorePemKeyPairCmd) *serviceresponse.StorePemKeyPairResponse
}
