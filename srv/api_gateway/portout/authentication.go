package portout

import (
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
)

// AuthenticationPortOut è una porta di output che si occupa di gestire l'autenticazione
type AuthenticationPortOut interface {
	// GetToken ritorna un token per l'utente fornito se le
	// credenziali sono valide, e TokenNone altrimenti.
	GetToken(username string) (types.UserToken, error)
	// GetRole ritorna il ruolo di un utente, dato il suo token.
	// Se il token non è valido, viene ritornato RoleNone
	GetRole(token types.UserToken) (types.UserRole, error)
}
