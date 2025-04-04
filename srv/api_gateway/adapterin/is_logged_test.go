package adapterin

import (
	"encoding/json"
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestIsLoggedOk(t *testing.T) {
	s := start(t)

	s.auth.EXPECT().ValidateToken("some.secure.jwt").Return(types.UserData{
		Username: "test",
		Role:     types.RoleGlobalAdmin,
	}, nil)

	req, _ := http.NewRequest("GET", s.base+"/api/v1/is_logged", nil)
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
	s := start(t)

	req, _ := http.NewRequest("GET", s.base+"/api/v1/is_logged", nil)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var respbody response.ResponseDTO[string]
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, "missing_token", respbody.Error)
}

func TestIsLoggedInvalidToken(t *testing.T) {
	s := start(t)

	s.auth.EXPECT().ValidateToken("some.secure.jwt").Return(types.UserData{}, business.ErrorTokenInvalid)

	req, _ := http.NewRequest("GET", s.base+"/api/v1/is_logged", nil)
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
	s := start(t)

	s.auth.EXPECT().ValidateToken("some.secure.jwt").Return(types.UserData{}, business.ErrorTokenExpired)

	req, _ := http.NewRequest("GET", s.base+"/api/v1/is_logged", nil)
	req.Header.Add("Authorization", "Bearer some.secure.jwt")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var respbody response.ResponseDTO[string]
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, "expired_token", respbody.Error)
}
