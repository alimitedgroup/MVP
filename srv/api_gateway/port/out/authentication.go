package out

type UserRole int

type UserToken string

const (
	LocalAdmin UserRole = iota
	GlobalAdmin
	Client
)

type AuthenticationPortOut interface {
	GetToken(username string) (UserToken, error)
	GetRole(token UserToken) (UserRole, error)
}
