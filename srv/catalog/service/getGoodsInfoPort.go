package service

type IGetGoodsInfoPort interface {
	GetGoodsInfo(ggqc *GetGoodsQuantityCmd) *GetGoodsInfoResponse
}
