package serviceportout

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/authenticator/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/authenticator/service/response"
)

type IStorePemKeyPairInterface interface {
	StorePemKeyPair(cmd *servicecmd.StorePemKeyPairCmd) *serviceresponse.StorePemKeyPairResponse
}
