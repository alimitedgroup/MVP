package adapterin

import (
	"encoding/json"
	"fmt"
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestLoginOk(t *testing.T) {
	s := start(t)

	s.auth.EXPECT().Login("user").Return(types.LoginResult{
		Token: "some.secure.token",
		Role:  types.RoleClient,
	}, nil)

	req, err := http.NewRequest(
		"POST",
		s.base+"/api/v1/login",
		strings.NewReader(url.Values{"username": {"user"}}.Encode()),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var respBody dto.AuthLoginResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)
	require.Equal(t, "some.secure.token", respBody.Token)
}

func TestLoginMissingUsername(t *testing.T) {
	s := start(t)

	req, err := http.NewRequest(
		"POST",
		s.base+"/api/v1/login",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var respBody response.ResponseDTO[dto.MissingRequiredFieldError]
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)
	require.Equal(t, "missing_field", respBody.Error)
	require.Equal(t, "username", respBody.Message.Field)
}

func TestLoginInternalError(t *testing.T) {
	s := start(t)

	s.auth.EXPECT().Login("user").Return(types.LoginResult{}, fmt.Errorf("some error"))

	req, err := http.NewRequest(
		"POST",
		s.base+"/api/v1/login",
		strings.NewReader(url.Values{"username": {"user"}}.Encode()),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var respBody response.ResponseDTO[string]
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)
	require.Equal(t, "internal_error", respBody.Error)
}

func TestLoginAuthFailed(t *testing.T) {
	s := start(t)

	s.auth.EXPECT().Login("user").Return(types.LoginResult{
		Token: "",
		Role:  types.RoleNone,
	}, nil)

	req, err := http.NewRequest(
		"POST",
		s.base+"/api/v1/login",
		strings.NewReader(url.Values{"username": {"user"}}.Encode()),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var respBody response.ResponseDTO[string]
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)
	require.Equal(t, "authentication_failed", respBody.Error)
}
