package service

type IAddOrChangeGoodDataPort interface {
	AddOrChangeGoodData(agc *AddGoodCmd) *AddOrChangeResponse
}
