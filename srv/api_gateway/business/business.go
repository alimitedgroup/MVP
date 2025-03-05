package business

import (
	"errors"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portout"
	"go.uber.org/fx"
	"time"
)

var (
	ErrorInvalidCredentials = errors.New("invalid credentials")
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
		return LoginResult{}, err
	}
	if token == "" {
		return LoginResult{}, ErrorInvalidCredentials
	}

	role, err := b.authAdapter.GetRole(token)
	if err != nil {
		return LoginResult{}, err
	}

	// TODO: bisognerebbe prendere la scadenza dall'output del servizio di Authentication
	expiration := time.Now().Add(7 * 24 * time.Hour)

	return LoginResult{
		Token:           token,
		TokenExpiration: expiration,
		Role:            role,
	}, nil
}

// LoginResult è il risultato di un login avvenuto con successo.
type LoginResult struct {
	// Token è una stringa opaca che il client dovrà fornire per autenticarsi.
	Token portout.UserToken
	// TokenExpiration è un tempo nel futuro in cui Token scadrà.
	// Quando ciò avviene, sarà necessario autenticarsi di nuovo.
	TokenExpiration time.Time
	// Role è il ruolo che è assegnato all'utente.
	Role portout.UserRole
}
