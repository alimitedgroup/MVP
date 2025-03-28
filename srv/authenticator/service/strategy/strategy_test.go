package serviceauthenticator

import (
	"testing"

	common "github.com/alimitedgroup/MVP/srv/authenticator/authCommon"
	serviceobject "github.com/alimitedgroup/MVP/srv/authenticator/service/object"
	"github.com/magiconair/properties/assert"
)

func TestUsername(t *testing.T) {
	sa := NewSimpleAuthAlg()
	udc := serviceobject.NewUserData("client")
	uda := serviceobject.NewUserData("global_admin")
	udla := serviceobject.NewUserData("local_admin")
	response, err := sa.Authenticate(*udc)
	assert.Equal(t, err, nil)
	assert.Equal(t, response, "client")
	response, err = sa.Authenticate(*uda)
	assert.Equal(t, err, nil)
	assert.Equal(t, response, "global_admin")
	response, err = sa.Authenticate(*udla)
	assert.Equal(t, err, nil)
	assert.Equal(t, response, "local_admin")
}

func TestWrongUsername(t *testing.T) {
	sa := NewSimpleAuthAlg()
	udc := serviceobject.NewUserData("test-username")
	response, err := sa.Authenticate(*udc)
	assert.Equal(t, err, common.ErrUserNotLegit)
	assert.Equal(t, response, "")
}
