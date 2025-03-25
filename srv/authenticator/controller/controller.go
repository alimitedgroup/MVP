package controller

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common"
	commonauth "github.com/alimitedgroup/MVP/srv/authenticator/authCommon"
	servicecmd "github.com/alimitedgroup/MVP/srv/authenticator/service/cmd"
	serviceportin "github.com/alimitedgroup/MVP/srv/authenticator/service/portIn"
	"github.com/nats-io/nats.go"
)

type authController struct {
	tokenUseCase serviceportin.IGetTokenUseCase
}

func NewAuthController(tokenUseCase serviceportin.IGetTokenUseCase) *authController {
	return &authController{tokenUseCase: tokenUseCase}
}

func (ar *authController) checkGetTokenRequest(dto *common.AuthLoginRequest) error {
	if dto.Username == "" {
		return commonauth.ErrUserNotLegit
	}
	return nil
}

func (ar *authController) NewTokenRequest(ctx context.Context, msg *nats.Msg) error {
	var dto common.AuthLoginRequest
	err := json.Unmarshal(msg.Data, &dto)
	if err != nil {
		response, err := json.Marshal(common.AuthLoginResponse{Token: ""})
		if err != nil {
			return err
		}
		err = msg.Respond(response)
		if err != nil {
			return err
		}
		return nil
	}
	err = ar.checkGetTokenRequest(&dto)

	if err != nil {
		response, err := json.Marshal(common.AuthLoginResponse{Token: ""})
		if err != nil {
			return err
		}
		err = msg.Respond(response)
		if err != nil {
			return err
		}
		return nil
	}

	tokenResponse := ar.tokenUseCase.GetToken(servicecmd.NewGetTokenCmd(dto.Username))

	if tokenResponse.GetError() != nil {
		response, err := json.Marshal(common.AuthLoginResponse{Token: ""})
		if err != nil {
			return err
		}
		err = msg.Respond(response)
		if err != nil {
			return err
		}
		return nil
	}

	response, err := json.Marshal(common.AuthLoginResponse{Token: tokenResponse.GetToken()})

	if err != nil {
		return err
	}

	err = msg.Respond(response)

	if err != nil {
		return err
	}

	return nil
}
