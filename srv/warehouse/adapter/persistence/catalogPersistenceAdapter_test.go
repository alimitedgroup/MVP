package persistence

import (
	"testing"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination mock_catalogRepository.go -package persistence github.com/alimitedgroup/MVP/srv/warehouse/adapter/persistence ICatalogRepository

func TestCatalogPersistanceAdapterSetAndGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockICatalogRepository(ctrl)

	mock.EXPECT().GetGood(gomock.Any()).Return(&Good{
		Id:          "1",
		Name:        "blue_hat",
		Description: "very beautiful hat",
	})
	mock.EXPECT().SetGood(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)

	a := NewCatalogPersistanceAdapter(mock)

	a.ApplyCatalogUpdate(model.GoodInfo{
		ID:          "1",
		Name:        "blue_hat",
		Description: "very beautiful hat",
	})

	goodInfo := a.GetGood("1")
	require.NotNil(t, goodInfo)
	require.Equal(t, goodInfo.ID, model.GoodID("1"))
}

func TestCatalogPersistanceAdapterGetNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockICatalogRepository(ctrl)

	mock.EXPECT().GetGood(gomock.Any()).Return(nil)

	a := NewCatalogPersistanceAdapter(mock)

	require.Nil(t, a.GetGood("1"))
}
