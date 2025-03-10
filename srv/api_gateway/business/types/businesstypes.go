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
