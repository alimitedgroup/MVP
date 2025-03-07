package types

type UserRole int

const (
	RoleLocalAdmin UserRole = iota
	RoleGlobalAdmin
	RoleClient
	RoleNone
)

type UserToken string

const TokenNone UserToken = ""
