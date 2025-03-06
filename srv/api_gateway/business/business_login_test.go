package business

import (
	"fmt"
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
		{"client", portout.RoleClient},
		{"local_admin", portout.RoleLocalAdmin},
		{"global_admin", portout.RoleGlobalAdmin},
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
	ctrl := gomock.NewController(t)
	mock := NewMockAuthenticationPortOut(ctrl)

	mock.EXPECT().GetToken(gomock.Any()).Return(portout.TokenNone, nil)

	business := NewBusiness(mock)
	_, err := business.Login("user")
	require.ErrorIs(t, err, ErrorInvalidCredentials)
}

func TestLoginGetTokenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockAuthenticationPortOut(ctrl)

	mock.EXPECT().GetToken(gomock.Any()).Return(portout.TokenNone, fmt.Errorf("some error"))

	business := NewBusiness(mock)
	_, err := business.Login("user")
	require.ErrorIs(t, err, ErrorGetToken)
}

func TestLoginGetRoleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockAuthenticationPortOut(ctrl)

	mock.EXPECT().GetToken(gomock.Any()).Return(portout.UserToken("some.secure.jwt"), nil)
	mock.EXPECT().GetRole(portout.UserToken("some.secure.jwt")).Return(portout.RoleNone, fmt.Errorf("some error"))

	business := NewBusiness(mock)
	_, err := business.Login("user")
	require.ErrorIs(t, err, ErrorGetRole)
}
