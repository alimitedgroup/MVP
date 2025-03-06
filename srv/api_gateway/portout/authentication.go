package portout

type UserRole int

const (
	RoleLocalAdmin UserRole = iota
	RoleGlobalAdmin
	RoleClient
	RoleNone
)

type UserToken string

const TokenNone UserToken = ""

// AuthenticationPortOut è una porta di output che si occupa di gestire l'autenticazione
type AuthenticationPortOut interface {
	// GetToken ritorna un token per l'utente fornito se le
	// credenziali sono valide, e TokenNone altrimenti.
	GetToken(username string) (UserToken, error)
	// GetRole ritorna il ruolo di un utente, dato il suo token.
	// Se il token non è valido, viene ritornato RoleNone
	GetRole(token UserToken) (UserRole, error)
}
