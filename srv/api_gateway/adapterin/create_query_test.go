package adapterin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateQuery(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	s.auth.EXPECT().ValidateToken("some.secure.jwt").Return(types.UserData{
		Username: "test",
		Role:     types.RoleGlobalAdmin,
	}, nil)
	s.notifications.EXPECT().CreateQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return("1", nil)

	payload := bytes.NewReader([]byte(`{"good_id": "1", "operator": ">", "threshold": 10}`))
	req, err := http.NewRequest(http.MethodPost, s.base+"/api/v1/notifications/queries", payload)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer some.secure.jwt")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var respbody dto.CreateQueryResponse
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, "1", respbody.QueryID)
}
