package controller

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/alimitedgroup/MVP/common"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	commonauth "github.com/alimitedgroup/MVP/srv/authenticator/authCommon"
	servicecmd "github.com/alimitedgroup/MVP/srv/authenticator/service/cmd"
	serviceportin "github.com/alimitedgroup/MVP/srv/authenticator/service/portIn"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	AuthRequests metric.Int64Counter
	Logger       *zap.Logger
	metricMap    sync.Map
)

type MetricParams struct {
	fx.In
	Logger *zap.Logger
	Meter  metric.Meter
}

type authController struct {
	tokenUseCase serviceportin.IGetTokenUseCase
}

func NewAuthController(tokenUseCase serviceportin.IGetTokenUseCase, mp MetricParams) *authController {
	observability.CounterSetup(&mp.Meter, mp.Logger, &AuthRequests, &metricMap, "num_token_requests")
	Logger = mp.Logger
	return &authController{tokenUseCase: tokenUseCase}
}

func (ar *authController) checkGetTokenRequest(dto *common.AuthLoginRequest) error {
	if dto.Username == "" {
		return commonauth.ErrUserNotLegit
	}
	return nil
}

func (ar *authController) NewTokenRequest(ctx context.Context, msg *nats.Msg) error {

	Logger.Info("Received new token generation request")
	verdict := "success"

	defer func() {
		Logger.Info("Generation token request terminated")
		AuthRequests.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	var dto common.AuthLoginRequest
	err := json.Unmarshal(msg.Data, &dto)
	if err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		err = broker.RespondToMsg(msg, common.AuthLoginResponse{Token: ""})
		if err != nil {
			Logger.Error("Cannot send response", zap.Error(err))
			return err
		}
		return nil
	}

	err = ar.checkGetTokenRequest(&dto)

	if err != nil {
		Logger.Debug("Bad username provided", zap.Error(err))
		err = broker.RespondToMsg(msg, common.AuthLoginResponse{Token: ""})
		if err != nil {
			Logger.Error("Cannot send response", zap.Error(err))
		}
		verdict = "bad request"
		return err
	}

	tokenResponse := ar.tokenUseCase.GetToken(servicecmd.NewGetTokenCmd(dto.Username))

	if tokenResponse.GetError() != nil {
		Logger.Debug("Cannot generate token", zap.Error(err))
		err = broker.RespondToMsg(msg, common.AuthLoginResponse{Token: ""})
		if err != nil {
			Logger.Error("Cannot send response", zap.Error(err))
		}
		verdict = "cannot generate token"
		return err
	}

	err = broker.RespondToMsg(msg, common.AuthLoginResponse{Token: tokenResponse.GetToken()})
	if err != nil {
		Logger.Error("Cannot send response", zap.Error(err))
		verdict = "token generated"
		return err
	}
	return nil
}
