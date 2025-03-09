package business

import (
	"errors"
	"fmt"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
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

func (b *Business) Login(username string) (LoginResult, error) {
	token, err := b.authAdapter.GetToken(username)
	if err != nil {
		return LoginResult{}, fmt.Errorf("%w: %w", ErrorGetToken, err)
	}
	if token == "" {
		return LoginResult{}, ErrorInvalidCredentials
	}

	parsed, err := b.authAdapter.VerifyToken(token)
	if err != nil {
		if errors.Is(err, portout.ErrTokenExpired) {
			return LoginResult{}, ErrorTokenExpired
		} else if errors.Is(err, portout.ErrTokenInvalid) {
			return LoginResult{}, ErrorTokenInvalid
		} else {
			return LoginResult{}, fmt.Errorf("%w: %w", ErrorGetToken, err)
		}
	}

	role, err := b.authAdapter.GetRole(parsed)
	if err != nil {
		return LoginResult{}, fmt.Errorf("%w: %w", ErrorGetRole, err)
	}

	// TODO: bisognerebbe prendere la scadenza dall'output del servizio di Authentication
	expiration := time.Now().Add(7 * 24 * time.Hour)

	return LoginResult{
		Token:           token,
		TokenExpiration: expiration,
		Role:            role,
	}, nil
}

func (b *Business) ValidateToken(token string) (UserData, error) {
	tok, err := b.authAdapter.VerifyToken(types.UserToken(token))
	if err != nil {
		if errors.Is(err, portout.ErrTokenExpired) {
			return UserData{}, ErrorTokenExpired
		} else if errors.Is(err, portout.ErrTokenInvalid) {
			return UserData{}, ErrorTokenInvalid
		} else {
			return UserData{}, fmt.Errorf("%w: %w", ErrorGetToken, err)
		}
	}

	username, err := b.authAdapter.GetUsername(tok)
	if err != nil {
		if errors.Is(err, portout.ErrTokenInvalid) {
			return UserData{}, ErrorTokenInvalid
		} else if errors.Is(err, portout.ErrTokenExpired) {
			return UserData{}, ErrorTokenExpired
		} else {
			return UserData{}, fmt.Errorf("%w: %w", ErrorGetUsername, err)
		}
	}

	role, err := b.authAdapter.GetRole(tok)
	if err != nil {
		if errors.Is(err, portout.ErrTokenInvalid) {
			return UserData{}, ErrorTokenInvalid
		} else if errors.Is(err, portout.ErrTokenExpired) {
			return UserData{}, ErrorTokenExpired
		} else {
			return UserData{}, fmt.Errorf("%w: %w", ErrorGetRole, err)
		}
	}

	return UserData{Username: username, Role: role}, err
}

// LoginResult è il risultato di un login avvenuto con successo.
type LoginResult struct {
	// Token è una stringa opaca che il client dovrà fornire per autenticarsi.
	Token types.UserToken
	// TokenExpiration è un tempo nel futuro in cui Token scadrà.
	// Quando ciò avviene, sarà necessario autenticarsi di nuovo.
	TokenExpiration time.Time
	// Role è il ruolo che è assegnato all'utente.
	Role types.UserRole
}

type UserData struct {
	Username string
	// Role è il ruolo che è assegnato all'utente.
	Role types.UserRole
}
