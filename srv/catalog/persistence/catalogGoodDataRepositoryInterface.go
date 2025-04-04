package persistence

import "github.com/alimitedgroup/MVP/common/dto"

type ICatalogGoodDataRepository interface {
	GetGoods() map[string]dto.Good
	AddGood(goodID string, name string, description string) error
}
