package business

import (
	"testing"
	"time"

	types "github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination mock_service.go -package business github.com/alimitedgroup/MVP/srv/notification/portout RuleRepository,StockEventPublisher,RuleQueryRepository,StockRepository

func TestAddQueryRule(t *testing.T) {
	ctrl := gomock.NewController(t)
	ruleRepoMock := NewMockRuleRepository(ctrl)
	alertPublisherMock := NewMockStockEventPublisher(ctrl)
	quantityReaderMock := NewMockRuleQueryRepository(ctrl)
	stockRepoMock := NewMockStockRepository(ctrl)

	ruleRepoMock.EXPECT().AddRule(gomock.Any()).Return(uuid.MustParse("391d2936-c37b-4294-bfdc-29e2473a5052"), nil)

	business := NewBusiness(ruleRepoMock, alertPublisherMock, quantityReaderMock, stockRepoMock)
	uuid, err := business.AddQueryRule(types.QueryRule{
		GoodId:    "1",
		Operator:  "<",
		Threshold: 10,
	})
	require.NoError(t, err)
	require.Equal(t, "391d2936-c37b-4294-bfdc-29e2473a5052", uuid.String())
}

func TestGetQueryRule(t *testing.T) {
	ctrl := gomock.NewController(t)
	ruleRepoMock := NewMockRuleRepository(ctrl)
	alertPublisherMock := NewMockStockEventPublisher(ctrl)
	quantityReaderMock := NewMockRuleQueryRepository(ctrl)
	stockRepoMock := NewMockStockRepository(ctrl)

	ruleRepoMock.EXPECT().GetRule(gomock.Any()).Return(types.QueryRule{
		GoodId:    "1",
		Operator:  "<",
		Threshold: 10,
	}, nil)

	business := NewBusiness(ruleRepoMock, alertPublisherMock, quantityReaderMock, stockRepoMock)
	rule, err := business.GetQueryRule(uuid.MustParse("391d2936-c37b-4294-bfdc-29e2473a5052"))
	require.NoError(t, err)
	require.Equal(t, "1", rule.GoodId)
}

func TestListRules(t *testing.T) {
	ctrl := gomock.NewController(t)
	ruleRepoMock := NewMockRuleRepository(ctrl)
	alertPublisherMock := NewMockStockEventPublisher(ctrl)
	quantityReaderMock := NewMockRuleQueryRepository(ctrl)
	stockRepoMock := NewMockStockRepository(ctrl)

	ruleRepoMock.EXPECT().ListRules().Return([]types.QueryRuleWithId{
		{
			RuleId: uuid.MustParse("391d2936-c37b-4294-bfdc-29e2473a5052"),
			QueryRule: types.QueryRule{
				GoodId:    "1",
				Operator:  "<",
				Threshold: 10,
			},
		},
	}, nil)

	business := NewBusiness(ruleRepoMock, alertPublisherMock, quantityReaderMock, stockRepoMock)
	rules, err := business.ListQueryRules()
	require.NoError(t, err)
	require.Len(t, rules, 1)
	require.Equal(t, "1", rules[0].GoodId)
	require.Equal(t, "391d2936-c37b-4294-bfdc-29e2473a5052", rules[0].RuleId.String())
}

func TestEditQueryRule(t *testing.T) {
	ctrl := gomock.NewController(t)
	ruleRepoMock := NewMockRuleRepository(ctrl)
	alertPublisherMock := NewMockStockEventPublisher(ctrl)
	quantityReaderMock := NewMockRuleQueryRepository(ctrl)
	stockRepoMock := NewMockStockRepository(ctrl)

	ruleRepoMock.EXPECT().EditRule(gomock.Any(), gomock.Any()).Return(nil)

	business := NewBusiness(ruleRepoMock, alertPublisherMock, quantityReaderMock, stockRepoMock)
	err := business.EditQueryRule(uuid.MustParse("391d2936-c37b-4294-bfdc-29e2473a5052"), types.EditRule{
		GoodId:    nil,
		Operator:  nil,
		Threshold: nil,
	})
	require.NoError(t, err)
}

func TestRemoveQueryRule(t *testing.T) {
	ctrl := gomock.NewController(t)
	ruleRepoMock := NewMockRuleRepository(ctrl)
	alertPublisherMock := NewMockStockEventPublisher(ctrl)
	quantityReaderMock := NewMockRuleQueryRepository(ctrl)
	stockRepoMock := NewMockStockRepository(ctrl)

	ruleRepoMock.EXPECT().RemoveRule(gomock.Any()).Return(nil)

	business := NewBusiness(ruleRepoMock, alertPublisherMock, quantityReaderMock, stockRepoMock)
	err := business.RemoveQueryRule(uuid.MustParse("391d2936-c37b-4294-bfdc-29e2473a5052"))
	require.NoError(t, err)
}

func TestPublishStockAlert(t *testing.T) {
	ctrl := gomock.NewController(t)
	ruleRepoMock := NewMockRuleRepository(ctrl)
	alertPublisherMock := NewMockStockEventPublisher(ctrl)
	quantityReaderMock := NewMockRuleQueryRepository(ctrl)
	stockRepoMock := NewMockStockRepository(ctrl)

	alertPublisherMock.EXPECT().PublishStockAlert(gomock.Any()).Return(nil)

	business := NewBusiness(ruleRepoMock, alertPublisherMock, quantityReaderMock, stockRepoMock)
	err := business.PublishStockAlert(types.StockAlertEvent{
		Id:              "391d2936-c37b-4294-bfdc-29e2473a5052",
		Status:          types.StockPending,
		GoodID:          "1",
		Operator:        "<",
		Threshold:       10,
		CurrentQuantity: 5,
		Timestamp:       time.Now().UnixMilli(),
	})
	require.NoError(t, err)
}

func TestGetCurrentQuantityByGoodID(t *testing.T) {
	ctrl := gomock.NewController(t)
	ruleRepoMock := NewMockRuleRepository(ctrl)
	alertPublisherMock := NewMockStockEventPublisher(ctrl)
	quantityReaderMock := NewMockRuleQueryRepository(ctrl)
	stockRepoMock := NewMockStockRepository(ctrl)

	quantityReaderMock.EXPECT().GetCurrentQuantityByGoodID(gomock.Any()).Return(&types.GetRuleResultResponse{
		GoodID:          "1",
		CurrentQuantity: 10,
		Err:             nil,
	})

	business := NewBusiness(ruleRepoMock, alertPublisherMock, quantityReaderMock, stockRepoMock)
	resp := business.GetCurrentQuantityByGoodID("1")
	require.NoError(t, resp.Err)
	require.Equal(t, "1", resp.GoodID)
}

func TestRecordStockUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	ruleRepoMock := NewMockRuleRepository(ctrl)
	alertPublisherMock := NewMockStockEventPublisher(ctrl)
	quantityReaderMock := NewMockRuleQueryRepository(ctrl)
	stockRepoMock := NewMockStockRepository(ctrl)

	stockRepoMock.EXPECT().SaveStockUpdate(gomock.Any()).Return(nil)

	business := NewBusiness(ruleRepoMock, alertPublisherMock, quantityReaderMock, stockRepoMock)
	err := business.RecordStockUpdate(&types.AddStockUpdateCmd{
		WarehouseID: "1",
		Type:        "add",
		OrderID:     "",
		TransferID:  "",
		Goods: []types.StockGood{
			{
				ID:       "1",
				Quantity: 10,
				Delta:    10,
			},
		},
		Timestamp: time.Now().UnixMilli(),
	})
	require.NoError(t, err)
}
