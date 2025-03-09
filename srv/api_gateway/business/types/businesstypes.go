package types

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
