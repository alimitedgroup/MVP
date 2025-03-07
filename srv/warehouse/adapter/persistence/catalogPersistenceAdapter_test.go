package persistence

import (
	"testing"

	"github.com/alimitedgroup/MVP/srv/warehouse/model"
	"github.com/magiconair/properties/assert"
	"go.uber.org/fx"
)

type catalogRepositoryMock struct {
	M map[string]Good
}

func NewCatalogRepositoryMock() *catalogRepositoryMock {
	return &catalogRepositoryMock{M: make(map[string]Good)}
}

func (s *catalogRepositoryMock) SetGood(goodId string, name string, description string) bool {
	s.M[goodId] = Good{
		Id:          goodId,
		Name:        name,
		Description: description,
	}
	return true
}

func (s *catalogRepositoryMock) GetGood(goodId string) *Good {
	good, exist := s.M[goodId]
	if !exist {
		return nil
	}

	return &good
}

func TestCatalogPersistanceAdapter(t *testing.T) {
	ctx := t.Context()

	mock := NewCatalogRepositoryMock()

	app := fx.New(
		fx.Supply(fx.Annotate(mock, fx.As(new(CatalogRepository)))),
		fx.Provide(NewCatalogPersistanceAdapter),
		fx.Invoke(func(a *CatalogPersistanceAdapter, stockRepo CatalogRepository) {
			good := model.GoodInfo{
				ID:          "1",
				Name:        "blue_hat",
				Description: "very beautiful hat",
			}

			err := a.ApplyCatalogUpdate(good)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, stockRepo.GetGood("1").Name, "blue_hat")
		}),
	)

	err := app.Start(ctx)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = app.Stop(ctx)
		if err != nil {
			t.Error(err)
		}
	}()
}
