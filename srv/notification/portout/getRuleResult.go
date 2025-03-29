package portout

import (
	serviceresponse "github.com/alimitedgroup/MVP/srv/notification/service/response"
)

type IRuleQueryRepository interface {
	GetCurrentQuantityByGoodID(goodID string) *serviceresponse.GetRuleResultResponse
}
