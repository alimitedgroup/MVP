package portin

import (
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"time"
)

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

type Auth interface {
	Login(username string) (LoginResult, error)
	ValidateToken(token string) (UserData, error)
}
