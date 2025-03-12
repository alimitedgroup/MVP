package serviceauthenticator

import serviceobject "github.com/alimitedgroup/MVP/srv/authenticator/service/object"

type SimpleAuthAlg struct {
	validRoles map[string]bool
}

func NewSimpleAuthAlg() *SimpleAuthAlg {
	roles := make(map[string]bool)
	roles["admin"] = true
	roles["local_admin"] = true
	roles["client"] = true
	return &SimpleAuthAlg{validRoles: roles}
}

func (saa *SimpleAuthAlg) Authenticate(us serviceobject.UserData) bool {
	_, presence := saa.validRoles[us.GetRole()]
	return presence
}
