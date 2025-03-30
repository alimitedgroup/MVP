package portin

import (
	"github.com/alimitedgroup/MVP/common/dto"
)

type Notifications interface {
	CreateQuery(goodId string, operator string, threshold int) (string, error)
	GetQueries() ([]dto.Query, error)
}
