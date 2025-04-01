package business

import (
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/srv/notification/config"
	"github.com/alimitedgroup/MVP/srv/notification/portin"
	"github.com/alimitedgroup/MVP/srv/notification/portout"
	types "github.com/alimitedgroup/MVP/srv/notification/types"
	uuid "github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination mock_queryrules.go -package business github.com/alimitedgroup/MVP/srv/notification/portin QueryRules

func TestRuleCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	quantityReaderMock := NewMockRuleQueryRepository(ctrl)
	alertPublisherMock := NewMockStockEventPublisher(ctrl)
	queryRulesMock := NewMockQueryRules(ctrl)

	queryRulesMock.EXPECT().ListQueryRules().Return([]types.QueryRuleWithId{
		{
			RuleId: uuid.MustParse("391d2936-c37b-4294-bfdc-29e2473a5052"),
			QueryRule: types.QueryRule{
				GoodId:    "1",
				Operator:  "<",
				Threshold: 10,
			},
		},
	}, nil).AnyTimes()
	quantityReaderMock.EXPECT().GetCurrentQuantityByGoodID(gomock.Any()).Return(&types.GetRuleResultResponse{
		GoodID:          "1",
		CurrentQuantity: 5,
		Err:             nil,
	}).AnyTimes()
	alertPublisherMock.EXPECT().PublishStockAlert(gomock.Any()).Return(nil).AnyTimes()
	alertPublisherMock.EXPECT().RevokeStockAlert(gomock.Any()).Return(nil).AnyTimes()

	module := fx.Options(
		observability.ModuleTest,
		fx.Supply(gomock.NewController(t), t),
		fx.Supply(fx.Annotate(quantityReaderMock, fx.As(new(portout.RuleQueryRepository)))),
		fx.Supply(fx.Annotate(alertPublisherMock, fx.As(new(portout.StockEventPublisher)))),
		fx.Supply(fx.Annotate(queryRulesMock, fx.As(new(portin.QueryRules)))),
	)
	app := fx.New(
		module,
		fx.Invoke(func(lc fx.Lifecycle, logger *zap.Logger, rules portin.QueryRules, queries portout.RuleQueryRepository, publish portout.StockEventPublisher) {
			cfg := config.NotificationConfig{
				CheckerTimer: 1 * time.Second,
				ServiceId:    "1",
			}
			checker := NewRuleChecker(lc, logger, rules, queries, publish, &cfg)
			go func() {
				time.Sleep(3 * time.Second)
				checker.stop <- true
			}()
		}),
	)
	err := app.Start(t.Context())
	require.NoError(t, err)

	defer func() {
		err := app.Stop(t.Context())
		require.NoError(t, err)
	}()
}
