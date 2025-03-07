package service_Response

import "github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"

type GetGoodsInfoResponse struct {
	goodMap map[string]catalogCommon.Good
}

func NewGetGoodsInfoResponse(goodMap map[string]catalogCommon.Good) *GetGoodsInfoResponse {
	return &GetGoodsInfoResponse{goodMap: goodMap}
}

func (ggqr *GetGoodsInfoResponse) GetMap() map[string]catalogCommon.Good {

	return ggqr.goodMap
}
