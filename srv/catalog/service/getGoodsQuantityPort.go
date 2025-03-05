package service

type IGetGoodsQuantityPort interface {
	GetGoodsQuantity(ggqc *GetGoodsQuantityCmd) *GetGoodsQuantityResponse
}
