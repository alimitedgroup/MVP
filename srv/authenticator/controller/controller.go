package controller

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/alimitedgroup/MVP/common"
	"github.com/alimitedgroup/MVP/common/lib/broker"
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
	Mutex        sync.Mutex
)

func setCounter(c metric.Int64Counter) {
	Mutex.Lock()
	defer Mutex.Unlock()
	AuthRequests = c
}

func incrementCounter(v string) {
	Mutex.Lock()
	defer Mutex.Unlock()
	AuthRequests.Add(context.Background(), 1, metric.WithAttributes(attribute.String("verdict", v)))
}

type MetricParams struct {
	fx.In
	Logger *zap.Logger
	Meter  metric.Meter
}

func counter(p MetricParams, name string, options ...metric.Int64CounterOption) metric.Int64Counter {
	ctr, err := (p.Meter).Int64Counter(name, options...)
	if err != nil {
		p.Logger.Fatal("Failed to setup OpenTelemetry counter", zap.String("name", name), zap.Error(err))
	}
	return ctr
}

type authController struct {
	tokenUseCase serviceportin.IGetTokenUseCase
}

func NewAuthController(tokenUseCase serviceportin.IGetTokenUseCase, mp MetricParams) *authController {
	setCounter(counter(mp, "num_token_requests"))
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
		//AuthRequests.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		incrementCounter(verdict)
	}()

	var dto common.AuthLoginRequest
	err := json.Unmarshal(msg.Data, &dto)
	if err != nil {
		Logger.Debug("Bad request", zap.Error(err))
		err = broker.RespondToMsg(msg, common.AuthLoginResponse{Token: ""})
		if err != nil {
			Logger.Error("Cannot send response", zap.Error(err))
			return nil
		}
		verdict = "bad request"
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
		return nil
	}

	tokenResponse := ar.tokenUseCase.GetToken(servicecmd.NewGetTokenCmd(dto.Username))

	if tokenResponse.GetError() != nil {
		Logger.Debug("Cannot generate token", zap.Error(err))
		err = broker.RespondToMsg(msg, common.AuthLoginResponse{Token: ""})
		if err != nil {
			Logger.Error("Cannot send response", zap.Error(err))
		}
		verdict = "cannot generate token"
		return nil
	}

	err = broker.RespondToMsg(msg, common.AuthLoginResponse{Token: tokenResponse.GetToken()})
	if err != nil {
		Logger.Error("Cannot send response", zap.Error(err))
		verdict = "token generated"
		return nil
	}
	Logger.Debug("Genereation token request terminated")
	return nil
}
