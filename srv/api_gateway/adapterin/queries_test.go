package adapterin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/stretchr/testify/require"
)

func TestGetQueries(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	s.auth.EXPECT().ValidateToken("some.secure.jwt").Return(types.UserData{
		Username: "test",
		Role:     types.RoleGlobalAdmin,
	}, nil)
	s.notifications.EXPECT().GetQueries().Return(
		[]dto.Query{
			{
				QueryID:   "b8f6208f-0828-469f-811d-748ffbfd24b6",
				GoodID:    "1",
				Operator:  ">",
				Threshold: 10,
			},
		},
		nil,
	)

	req, err := http.NewRequest(http.MethodGet, s.base+"/api/v1/notifications/queries", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer some.secure.jwt")
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var respbody dto.GetQueriesResponse
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, []dto.Query{
		{
			QueryID:   "b8f6208f-0828-469f-811d-748ffbfd24b6",
			GoodID:    "1",
			Operator:  ">",
			Threshold: 10,
		},
	}, respbody.Queries)
}

func TestGetQueriesError(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	s.auth.EXPECT().ValidateToken("some.secure.jwt").Return(types.UserData{
		Username: "test",
		Role:     types.RoleGlobalAdmin,
	}, nil)

	s.notifications.EXPECT().GetQueries().Return([]dto.Query{}, fmt.Errorf("some error"))

	req, err := http.NewRequest(http.MethodGet, s.base+"/api/v1/notifications/queries", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer some.secure.jwt")
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var respbody response.ResponseDTO[string]
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, "internal_error", respbody.Error)
}
