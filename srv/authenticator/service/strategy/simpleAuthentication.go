package serviceauthenticator

import (
	common "github.com/alimitedgroup/MVP/srv/authenticator/authCommon"
	serviceobject "github.com/alimitedgroup/MVP/srv/authenticator/service/object"
)

type SimpleAuthAlg struct {
	usernameRoles map[string]string
}

func NewSimpleAuthAlg() *SimpleAuthAlg {
	roles := make(map[string]string)
	roles["admin"] = "admin"
	roles["local_admin"] = "local_admin"
	roles["client"] = "client"
	return &SimpleAuthAlg{usernameRoles: roles}
}

func (saa *SimpleAuthAlg) Authenticate(us serviceobject.UserData) (string, error) {
	role, presence := saa.usernameRoles[us.GetUsername()]
	if presence {
		return role, nil
	}
	return "", common.ErrUserNotLegit
}
