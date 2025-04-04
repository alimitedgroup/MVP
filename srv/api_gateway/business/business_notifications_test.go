package business

import (
	"testing"

	"github.com/alimitedgroup/MVP/common/dto"
	types "github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination mock_notifications.go -package business github.com/alimitedgroup/MVP/srv/api_gateway/portout NotificationPortOut

func TestGetQueries(t *testing.T) {
	ctrl := gomock.NewController(t)
	authMock := NewMockAuthenticationPortOut(ctrl)
	catalogMock := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)
	notificationsMock := NewMockNotificationPortOut(ctrl)

	notificationsMock.EXPECT().GetQueries().Return([]types.QueryRuleWithId{
		{
			RuleId: uuid.MustParse("4ab2cff3-9cc9-4b2e-b982-aa467565e676"),
			QueryRule: types.QueryRule{
				GoodId:    "1",
				Operator:  ">",
				Threshold: 10,
			},
		},
		{
			RuleId: uuid.MustParse("b8f6208f-0828-469f-811d-748ffbfd24b6"),
			QueryRule: types.QueryRule{
				GoodId:    "2",
				Operator:  ">",
				Threshold: 10,
			},
		},
	}, nil)

	p := BusinessParams{
		Auth:         authMock,
		Catalog:      catalogMock,
		Order:        orderMock,
		Notification: notificationsMock,
		Logger:       zaptest.NewLogger(t),
	}
	business := NewBusiness(p)
	queries, err := business.GetQueries()
	require.NoError(t, err)
	require.Len(t, queries, 2)
	require.ElementsMatch(t, []dto.Query{
		{
			QueryID:   "4ab2cff3-9cc9-4b2e-b982-aa467565e676",
			GoodID:    "1",
			Operator:  ">",
			Threshold: 10,
		},
		{
			QueryID:   "b8f6208f-0828-469f-811d-748ffbfd24b6",
			GoodID:    "2",
			Operator:  ">",
			Threshold: 10,
		},
	}, queries)
}

func TestCreateQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	authMock := NewMockAuthenticationPortOut(ctrl)
	catalogMock := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)
	notificationsMock := NewMockNotificationPortOut(ctrl)

	notificationsMock.EXPECT().CreateQuery(gomock.Any()).Return("1", nil)

	p := BusinessParams{
		Auth:         authMock,
		Catalog:      catalogMock,
		Order:        orderMock,
		Notification: notificationsMock,
		Logger:       zaptest.NewLogger(t),
	}
	business := NewBusiness(p)
	queryId, err := business.CreateQuery("1", ">", 10)
	require.NoError(t, err)
	require.Equal(t, "1", queryId)
}
