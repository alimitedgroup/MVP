package business

import (
	"errors"
	"fmt"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"time"

	"github.com/alimitedgroup/MVP/srv/api_gateway/portout"
	"go.uber.org/fx"
)

var (
	ErrorGetToken           = errors.New("error getting token for given credentials")
	ErrorGetRole            = errors.New("error getting role for given token")
	ErrorGetUsername        = errors.New("error getting username")
	ErrorInvalidCredentials = errors.New("invalid credentials")
	ErrorTokenInvalid       = errors.New("this token is invalid")
	ErrorTokenExpired       = errors.New("this token is expired")
)

type Business struct {
	authAdapter portout.AuthenticationPortOut
}

func NewBusiness(authAdapter portout.AuthenticationPortOut) *Business {
	return &Business{authAdapter: authAdapter}
}

var Module = fx.Options(
	fx.Provide(NewBusiness),
)

func (b *Business) Login(username string) (portin.LoginResult, error) {
	token, err := b.authAdapter.GetToken(username)
	if err != nil {
		return portin.LoginResult{}, fmt.Errorf("%w: %w", ErrorGetToken, err)
	}
	if token == "" {
		return portin.LoginResult{}, ErrorInvalidCredentials
	}

	parsed, err := b.authAdapter.VerifyToken(token)
	if err != nil {
		return portin.LoginResult{}, fmt.Errorf("%w: %w", ErrorGetToken, err)
	}

	role, err := b.authAdapter.GetRole(parsed)
	if err != nil {
		return portin.LoginResult{}, fmt.Errorf("%w: %w", ErrorGetRole, err)
	}

	// TODO: bisognerebbe prendere la scadenza dall'output del servizio di Authentication
	expiration := time.Now().Add(7 * 24 * time.Hour)

	return portin.LoginResult{
		Token:           token,
		TokenExpiration: expiration,
		Role:            role,
	}, nil
}

func (b *Business) ValidateToken(token string) (portin.UserData, error) {
	tok, err := b.authAdapter.VerifyToken(types.UserToken(token))
	if err != nil {
		if errors.Is(err, portout.ErrTokenExpired) {
			return portin.UserData{}, ErrorTokenExpired
		} else if errors.Is(err, portout.ErrTokenInvalid) {
			return portin.UserData{}, ErrorTokenInvalid
		} else {
			return portin.UserData{}, fmt.Errorf("%w: %w", ErrorGetToken, err)
		}
	}

	username, err := b.authAdapter.GetUsername(tok)
	if err != nil {
		if errors.Is(err, portout.ErrTokenInvalid) {
			return portin.UserData{}, ErrorTokenInvalid
		} else if errors.Is(err, portout.ErrTokenExpired) {
			return portin.UserData{}, ErrorTokenExpired
		} else {
			return portin.UserData{}, fmt.Errorf("%w: %w", ErrorGetUsername, err)
		}
	}

	role, err := b.authAdapter.GetRole(tok)
	if err != nil {
		if errors.Is(err, portout.ErrTokenInvalid) {
			return portin.UserData{}, ErrorTokenInvalid
		} else if errors.Is(err, portout.ErrTokenExpired) {
			return portin.UserData{}, ErrorTokenExpired
		} else {
			return portin.UserData{}, fmt.Errorf("%w: %w", ErrorGetRole, err)
		}
	}

	return portin.UserData{Username: username, Role: role}, err
}

// Asserzione a compile time che Business implementi le interfaccie delle porte di input
var _ portin.Auth = (*Business)(nil)
