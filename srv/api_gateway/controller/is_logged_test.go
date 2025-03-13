package controller

import (
	"encoding/json"
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestIsLoggedOk(t *testing.T) {
	auth, base := start(t)

	auth.EXPECT().ValidateToken("some.secure.jwt").Return(portin.UserData{
		Username: "test",
		Role:     types.RoleGlobalAdmin,
	}, nil)

	req, _ := http.NewRequest("GET", base+"/api/v1/is_logged", nil)
	req.Header.Add("Authorization", "Bearer some.secure.jwt")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var respbody dto.IsLoggedResponse
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, "global_admin", respbody.Role)
}

func TestIsLoggedMissingToken(t *testing.T) {
	_, base := start(t)

	req, _ := http.NewRequest("GET", base+"/api/v1/is_logged", nil)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var respbody response.ResponseDTO[string]
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, "missing_token", respbody.Error)
}

func TestIsLoggedInvalidToken(t *testing.T) {
	auth, base := start(t)

	auth.EXPECT().ValidateToken("some.secure.jwt").Return(portin.UserData{}, business.ErrorTokenInvalid)

	req, _ := http.NewRequest("GET", base+"/api/v1/is_logged", nil)
	req.Header.Add("Authorization", "Bearer some.secure.jwt")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var respbody response.ResponseDTO[string]
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, "invalid_token", respbody.Error)
}

func TestIsLoggedExpiredToken(t *testing.T) {
	auth, base := start(t)

	auth.EXPECT().ValidateToken("some.secure.jwt").Return(portin.UserData{}, business.ErrorTokenExpired)

	req, _ := http.NewRequest("GET", base+"/api/v1/is_logged", nil)
	req.Header.Add("Authorization", "Bearer some.secure.jwt")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var respbody response.ResponseDTO[string]
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, "expired_token", respbody.Error)
}
