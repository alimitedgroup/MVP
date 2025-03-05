package business

import (
	"github.com/alimitedgroup/MVP/srv/api_gateway/portout"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination auth_mock.go -package business github.com/alimitedgroup/MVP/srv/api_gateway/portout AuthenticationPortOut

func TestLogin(t *testing.T) {
	cases := []struct {
		string
		portout.UserRole
	}{
		{"client", portout.Client},
		{"local_admin", portout.LocalAdmin},
		{"global_admin", portout.GlobalAdmin},
	}

	for _, c := range cases {
		t.Run(c.string, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := NewMockAuthenticationPortOut(ctrl)

			mock.EXPECT().GetToken(c.string).Return(portout.UserToken("some.secure.jwt"), nil)
			mock.EXPECT().GetRole(portout.UserToken("some.secure.jwt")).Return(c.UserRole, nil)

			business := NewBusiness(mock)
			result, err := business.Login(c.string)
			require.NoError(t, err)
			require.Equal(t, portout.UserToken("some.secure.jwt"), result.Token)
			require.Equal(t, c.UserRole, result.Role)
			require.True(t, time.Now().Add(6*24*time.Hour).Before(result.TokenExpiration))
			require.True(t, time.Now().Add(8*24*time.Hour).After(result.TokenExpiration))
		})
	}
}

func TestLoginNoSuchUser(t *testing.T) {
	// TODO
}

func TestLoginGetTokenError(t *testing.T) {
	// TODO
}

func TestLoginGetRoleError(t *testing.T) {
	// TODO
}
