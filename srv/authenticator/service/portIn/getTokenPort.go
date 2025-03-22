package serviceportin

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/authenticator/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/authenticator/service/response"
)

type IGetTokenUseCase interface {
	GetToken(cmd *servicecmd.GetTokenCmd) *serviceresponse.GetTokenResponse
}
