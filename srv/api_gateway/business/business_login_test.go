package business

import (
	"fmt"
	"testing"

	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portout"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination mock_auth.go -package business github.com/alimitedgroup/MVP/srv/api_gateway/portout AuthenticationPortOut

func TestLogin(t *testing.T) {
	var token types.ParsedToken = struct{ test int }{}

	cases := []struct {
		string
		types.UserRole
	}{
		{"client", types.RoleClient},
		{"local_admin", types.RoleLocalAdmin},
		{"global_admin", types.RoleGlobalAdmin},
	}

	for _, c := range cases {
		t.Run(c.string, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			authMock := NewMockAuthenticationPortOut(ctrl)
			catalogMock := NewMockCatalogPortOut(ctrl)
			orderMock := NewMockOrderPortOut(ctrl)
			notificationsMock := NewMockNotificationPortOut(ctrl)

			authMock.EXPECT().GetToken(gomock.Any()).Return(types.UserToken("some.secure.jwt"), nil)
			authMock.EXPECT().VerifyToken(types.UserToken("some.secure.jwt")).Return(token, nil)
			authMock.EXPECT().GetRole(token).Return(c.UserRole, nil)

			p := BusinessParams{
				Auth:         authMock,
				Catalog:      catalogMock,
				Order:        orderMock,
				Notification: notificationsMock,
				Logger:       zaptest.NewLogger(t),
			}
			business := NewBusiness(p)
			result, err := business.Login(c.string)
			require.NoError(t, err)
			require.Equal(t, types.UserToken("some.secure.jwt"), result.Token)
			require.Equal(t, c.UserRole, result.Role)
		})
	}
}

func TestLoginNoSuchUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	authMock := NewMockAuthenticationPortOut(ctrl)
	catalogMock := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)
	notificationsMock := NewMockNotificationPortOut(ctrl)

	authMock.EXPECT().GetToken(gomock.Any()).Return(types.TokenNone, nil)

	p := BusinessParams{
		Auth:         authMock,
		Catalog:      catalogMock,
		Order:        orderMock,
		Notification: notificationsMock,
		Logger:       zaptest.NewLogger(t),
	}
	business := NewBusiness(p)
	_, err := business.Login("user")
	require.ErrorIs(t, err, ErrorInvalidCredentials)
}

func TestLoginGetTokenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	authMock := NewMockAuthenticationPortOut(ctrl)
	catalogMock := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)
	notificationsMock := NewMockNotificationPortOut(ctrl)

	authMock.EXPECT().GetToken(gomock.Any()).Return(types.TokenNone, fmt.Errorf("some error"))

	p := BusinessParams{
		Auth:         authMock,
		Catalog:      catalogMock,
		Order:        orderMock,
		Notification: notificationsMock,
		Logger:       zaptest.NewLogger(t),
	}
	business := NewBusiness(p)
	_, err := business.Login("user")
	require.ErrorIs(t, err, ErrorGetToken)
}

func TestLoginGetRoleError(t *testing.T) {
	var token types.ParsedToken = struct{ test int }{}

	ctrl := gomock.NewController(t)
	authMock := NewMockAuthenticationPortOut(ctrl)
	catalogMock := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)
	notificationsMock := NewMockNotificationPortOut(ctrl)

	authMock.EXPECT().GetToken(gomock.Any()).Return(types.UserToken("some.secure.jwt"), nil)
	authMock.EXPECT().VerifyToken(types.UserToken("some.secure.jwt")).Return(token, nil)
	authMock.EXPECT().GetRole(token).Return(types.RoleNone, fmt.Errorf("some error"))

	p := BusinessParams{
		Auth:         authMock,
		Catalog:      catalogMock,
		Order:        orderMock,
		Notification: notificationsMock,
		Logger:       zaptest.NewLogger(t),
	}
	business := NewBusiness(p)
	_, err := business.Login("user")
	require.ErrorIs(t, err, ErrorGetRole)
}

func TestLoginVerifyTokenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	authMock := NewMockAuthenticationPortOut(ctrl)
	catalogMock := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)
	notificationsMock := NewMockNotificationPortOut(ctrl)

	authMock.EXPECT().GetToken(gomock.Any()).Return(types.UserToken("some.secure.jwt"), nil)
	authMock.EXPECT().VerifyToken(types.UserToken("some.secure.jwt")).Return(nil, fmt.Errorf("some error"))

	p := BusinessParams{
		Auth:         authMock,
		Catalog:      catalogMock,
		Order:        orderMock,
		Notification: notificationsMock,
		Logger:       zaptest.NewLogger(t),
	}
	business := NewBusiness(p)
	_, err := business.Login("user")
	require.ErrorIs(t, err, ErrorGetToken)
}

func TestVerifyToken(t *testing.T) {
	var token types.ParsedToken = struct{ test int }{}

	ctrl := gomock.NewController(t)
	authMock := NewMockAuthenticationPortOut(ctrl)
	catalogMock := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)
	notificationsMock := NewMockNotificationPortOut(ctrl)

	authMock.EXPECT().VerifyToken(types.UserToken("some.secure.jwt")).Return(token, nil)
	authMock.EXPECT().GetUsername(token).Return("admin", nil)
	authMock.EXPECT().GetRole(token).Return(types.RoleClient, nil)

	p := BusinessParams{
		Auth:         authMock,
		Catalog:      catalogMock,
		Order:        orderMock,
		Notification: notificationsMock,
		Logger:       zaptest.NewLogger(t),
	}
	business := NewBusiness(p)
	data, err := business.ValidateToken("some.secure.jwt")
	require.NoError(t, err)
	require.Equal(t, data.Username, "admin")
	require.Equal(t, data.Role, types.RoleClient)
}

func TestVerifyTokenErrors(t *testing.T) {
	type testcase struct {
		port        error
		business    error
		description string
	}

	cases := []testcase{
		{portout.ErrTokenExpired, ErrorTokenExpired, "TokenExpired"},
		{portout.ErrTokenInvalid, ErrorTokenInvalid, "TokenInvalid"},
		{fmt.Errorf("some error"), ErrorGetToken, "GetToken"},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			authMock := NewMockAuthenticationPortOut(ctrl)
			catalogMock := NewMockCatalogPortOut(ctrl)
			orderMock := NewMockOrderPortOut(ctrl)
			notificationsMock := NewMockNotificationPortOut(ctrl)

			authMock.EXPECT().VerifyToken(types.UserToken("some.secure.jwt")).Return(42, c.port)

			p := BusinessParams{
				Auth:         authMock,
				Catalog:      catalogMock,
				Order:        orderMock,
				Notification: notificationsMock,
				Logger:       zaptest.NewLogger(t),
			}
			business := NewBusiness(p)
			data, err := business.ValidateToken("some.secure.jwt")
			require.ErrorIs(t, err, c.business)
			require.Zero(t, data)
		})
	}

}

func TestVerifyUsernameError(t *testing.T) {
	var token types.ParsedToken = struct{ test int }{}
	type testcase struct {
		port        error
		business    error
		description string
	}

	cases := []testcase{
		{portout.ErrTokenExpired, ErrorTokenExpired, "TokenExpired"},
		{portout.ErrTokenInvalid, ErrorTokenInvalid, "TokenInvalid"},
		{fmt.Errorf("some error"), ErrorGetUsername, "GetUsername"},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			authMock := NewMockAuthenticationPortOut(ctrl)
			catalogMock := NewMockCatalogPortOut(ctrl)
			orderMock := NewMockOrderPortOut(ctrl)
			notificationsMock := NewMockNotificationPortOut(ctrl)

			authMock.EXPECT().VerifyToken(types.UserToken("some.secure.jwt")).Return(token, nil)
			authMock.EXPECT().GetUsername(token).Return("", c.port)

			p := BusinessParams{
				Auth:         authMock,
				Catalog:      catalogMock,
				Order:        orderMock,
				Notification: notificationsMock,
				Logger:       zaptest.NewLogger(t),
			}
			business := NewBusiness(p)
			data, err := business.ValidateToken("some.secure.jwt")

			require.Zero(t, data)
			require.ErrorIs(t, err, c.business)
		})
	}
}

func TestVerifyRoleError(t *testing.T) {
	var token types.ParsedToken = struct{ test int }{}
	type testcase struct {
		port        error
		business    error
		description string
	}

	cases := []testcase{
		{portout.ErrTokenExpired, ErrorTokenExpired, "TokenExpired"},
		{portout.ErrTokenInvalid, ErrorTokenInvalid, "TokenInvalid"},
		{fmt.Errorf("some error"), ErrorGetRole, "GetRole"},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			authMock := NewMockAuthenticationPortOut(ctrl)
			catalogMock := NewMockCatalogPortOut(ctrl)
			orderMock := NewMockOrderPortOut(ctrl)
			notificationsMock := NewMockNotificationPortOut(ctrl)

			authMock.EXPECT().VerifyToken(types.UserToken("some.secure.jwt")).Return(token, nil)
			authMock.EXPECT().GetUsername(token).Return("admin", nil)
			authMock.EXPECT().GetRole(token).Return(types.RoleNone, c.port)

			p := BusinessParams{
				Auth:         authMock,
				Catalog:      catalogMock,
				Order:        orderMock,
				Notification: notificationsMock,
				Logger:       zaptest.NewLogger(t),
			}
			business := NewBusiness(p)
			data, err := business.ValidateToken("some.secure.jwt")

			require.Zero(t, data)
			require.ErrorIs(t, err, c.business)
		})
	}
}
