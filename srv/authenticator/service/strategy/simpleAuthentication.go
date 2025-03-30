package serviceauthenticator

import (
	"sync"

	common "github.com/alimitedgroup/MVP/srv/authenticator/authCommon"
	serviceobject "github.com/alimitedgroup/MVP/srv/authenticator/service/object"
)

type SimpleAuthAlg struct {
	usernameRoles map[string]string
	mutex         sync.Mutex
}

func NewSimpleAuthAlg() *SimpleAuthAlg {
	roles := make(map[string]string)
	roles["global_admin"] = "global_admin"
	roles["local_admin"] = "local_admin"
	roles["client"] = "client"
	return &SimpleAuthAlg{usernameRoles: roles}
}

func (saa *SimpleAuthAlg) Authenticate(us serviceobject.UserData) (string, error) {
	saa.mutex.Lock()
	defer saa.mutex.Unlock()
	role, presence := saa.usernameRoles[us.GetUsername()]
	if presence {
		return role, nil
	}
	return "", common.ErrUserNotLegit
}
