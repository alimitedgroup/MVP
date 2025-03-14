package serviceresponse

import (
	"github.com/alimitedgroup/MVP/common/dto"
)

type GetGoodsInfoResponse struct {
	goodMap map[string]dto.Good
}

func NewGetGoodsInfoResponse(goodMap map[string]dto.Good) *GetGoodsInfoResponse {
	return &GetGoodsInfoResponse{goodMap: goodMap}
}

func (ggqr *GetGoodsInfoResponse) GetMap() map[string]dto.Good {

	return ggqr.goodMap
}
