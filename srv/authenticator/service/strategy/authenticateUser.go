package serviceauthenticator

import serviceobject "github.com/alimitedgroup/MVP/srv/authenticator/service/object"

type IAuthenticateUserStrategy interface {
	Authenticate(us serviceobject.UserData) (string, error)
}
