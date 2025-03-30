package portout

import (
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/notification/types"
)

type NotificationPortOut interface {
	CreateQuery(dto dto.Rule) (string, error)
	GetQueries() ([]types.QueryRuleWithId, error)
}
