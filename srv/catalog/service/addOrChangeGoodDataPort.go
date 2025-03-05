package service

type IAddOrChangeGoodDataPort interface {
	AddOrChangeGoodData(agc *AddChangeGoodCmd) *AddOrChangeResponse
}
