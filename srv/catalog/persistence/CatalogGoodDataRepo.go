package persistence

import (
	"sync"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"
)

type CatalogGoodDataRepository struct {
	goodMap map[string]*dto.Good
	mutex   sync.Mutex
}

func NewCatalogGoodDataRepository() *CatalogGoodDataRepository {
	return &CatalogGoodDataRepository{goodMap: make(map[string]*dto.Good)}
}

func (cr *CatalogGoodDataRepository) GetGoods() map[string]dto.Good {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	result := make(map[string]dto.Good)
	for key := range cr.goodMap {
		result[key] = *cr.goodMap[key]
	}
	return result
}

func (cr *CatalogGoodDataRepository) AddGood(goodID string, name string, description string) error {
	cr.mutex.Lock()
	_, presence := cr.goodMap[goodID]
	if presence {
		cr.mutex.Unlock()
		return cr.changeGoodData(goodID, name, description)
	}
	cr.goodMap[goodID] = dto.NewGood(goodID, name, description)
	cr.mutex.Unlock()
	return nil
}

func (cr *CatalogGoodDataRepository) changeGoodData(goodID string, newName string, newDescription string) error {
	cr.mutex.Lock()
	_, presence := cr.goodMap[goodID]
	if !presence {
		cr.mutex.Unlock()
		return catalogCommon.ErrGoodIdNotValid
	}
	err := cr.goodMap[goodID].SetName(newName)
	if err != nil {
		cr.mutex.Unlock()
		return err
	}
	err = cr.goodMap[goodID].SetDescription(newDescription)
	if err != nil {
		cr.mutex.Unlock()
		return err
	}
	cr.mutex.Unlock()
	return nil
}
