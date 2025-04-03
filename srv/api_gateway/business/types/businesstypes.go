package types

import "fmt"

type UserRole int

const (
	RoleNone       UserRole = 0
	RoleLocalAdmin UserRole = iota
	RoleGlobalAdmin
	RoleClient
)

type UserToken string

const TokenNone UserToken = ""

type ParsedToken interface{}

func (role UserRole) String() string {
	switch role {
	case RoleLocalAdmin:
		return "local_admin"
	case RoleGlobalAdmin:
		return "global_admin"
	case RoleClient:
		return "client"
	default:
		panic(fmt.Sprintf("unknown role %d", role))
	}
}

func RoleFromString(s string) UserRole {
	switch s {
	case "client":
		return RoleClient
	case "global_admin":
		return RoleGlobalAdmin
	case "local_admin":
		return RoleLocalAdmin
	default:
		return RoleNone
	}
}

// LoginResult è il risultato di un login avvenuto con successo.
type LoginResult struct {
	// Token è una stringa opaca che il client dovrà fornire per autenticarsi.
	Token UserToken
	// Role è il ruolo che è assegnato all'utente.
	Role UserRole
}

type UserData struct {
	Username string
	// Role è il ruolo che è assegnato all'utente.
	Role UserRole
}

type WarehouseOverview struct {
	ID string
}
