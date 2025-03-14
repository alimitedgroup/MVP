package portout

import (
	"errors"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
)

var (
	// ErrTokenExpired viene ritornato se il token fornito era scaduto
	ErrTokenExpired = errors.New("token expired")
	// ErrTokenInvalid viene ritornato se il token fornito era invalido, ma non scaduto
	ErrTokenInvalid = errors.New("token is invalid")
)

// AuthenticationPortOut è una porta di output che si occupa di gestire l'autenticazione
type AuthenticationPortOut interface {
	// GetToken ritorna un token per l'utente fornito se le
	// credenziali sono valide, e TokenNone altrimenti.
	GetToken(username string) (types.UserToken, error)
	// GetUsername ritorna l'username un utente, dato il suo token.
	// Se il token non è valido, viene ritornata la stringa vuota
	GetUsername(token types.ParsedToken) (string, error)
	// GetRole ritorna il ruolo di un utente, dato il suo token.
	// Se il token non è valido, viene ritornato RoleNone.
	GetRole(token types.ParsedToken) (types.UserRole, error)
	// VerifyToken verifica se il token fornito è valido, non ritornando nessun errore se lo è.
	// Se il token è invalido, viene ritornato ErrTokenExpired se era scaduto, ErrTokenInvalid altrimenti
	VerifyToken(token types.UserToken) (types.ParsedToken, error)
}
