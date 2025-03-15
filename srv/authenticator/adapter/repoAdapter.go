package adapter

import (
	common "github.com/alimitedgroup/MVP/srv/authenticator/authCommon"
	"github.com/alimitedgroup/MVP/srv/authenticator/persistence"
	servicecmd "github.com/alimitedgroup/MVP/srv/authenticator/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/authenticator/service/response"
)

type AuthAdapter struct {
	repo persistence.IAuthPersistance
}

func NewAuthAdapter(repo persistence.IAuthPersistance) *AuthAdapter {
	return &AuthAdapter{repo: repo}
}

func (aa *AuthAdapter) StorePemKeyPair(cmd *servicecmd.StorePemKeyPairCmd) *serviceresponse.StorePemKeyPairResponse {
	if cmd.GetPemPrivateKey() == nil || cmd.GetPemPublicKey() == nil {
		return serviceresponse.NewStorePemKeyPairResponse(common.ErrKeyPairNotValid)
	}
	return serviceresponse.NewStorePemKeyPairResponse(aa.repo.StorePemKeyPair(*cmd.GetPemPrivateKey(), *cmd.GetPemPublicKey(), cmd.GetIssuer()))
}

func (aa *AuthAdapter) GetPemPrivateKey(cmd *servicecmd.GetPemPrivateKeyCmd) *serviceresponse.GetPemPrivateKeyResponse {
	prk, err := aa.repo.GetPemPrivateKey()
	if err != nil {
		return serviceresponse.NewGetPemPrivateKeyResponse(nil, "", err)
	}
	data := prk.GetBytes()
	return serviceresponse.NewGetPemPrivateKeyResponse(&data, prk.GetIssuer(), nil)
}

func (aa *AuthAdapter) GetPemPublicKey(cmd *servicecmd.GetPemPublicKeyCmd) *serviceresponse.GetPemPublicKeyResponse {
	puk, err := aa.repo.GetPemPublicKey()
	if err != nil {
		return serviceresponse.NewGetPemPublicKeyResponse(nil, "", err)
	}
	data := puk.GetBytes()
	return serviceresponse.NewGetPemPublicKeyResponse(&data, puk.GetIssuer(), nil)
}

func (aa *AuthAdapter) CheckKeyPairExistance(cmd *servicecmd.CheckPemKeyPairExistenceCmd) *serviceresponse.CheckKeyPairExistenceResponse {
	return serviceresponse.NewCheckKeyPairExistenceResponse(aa.repo.CheckKeyPairExistence())
}
