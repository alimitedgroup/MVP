package service

type ISetGoodQuantityPort interface {
	SetGoodQuantity(agqc *SetGoodQuantityCmd) *SetGoodQuantityResponse
}
