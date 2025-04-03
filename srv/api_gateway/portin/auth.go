package portin

import (
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
)

type Auth interface {
	Login(username string) (types.LoginResult, error)
	ValidateToken(token string) (types.UserData, error)
}
