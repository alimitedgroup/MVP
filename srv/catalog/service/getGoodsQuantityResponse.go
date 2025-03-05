package service

type GetGoodsQuantityResponse struct {
	goodMap map[string]int64
}

func NewGetGoodsQuantityResponse(goodMap map[string]int64) *GetGoodsQuantityResponse {
	return &GetGoodsQuantityResponse{goodMap: goodMap}
}

func (ggqr *GetGoodsQuantityResponse) GetMap() map[string]int64 {
	return ggqr.goodMap
}
