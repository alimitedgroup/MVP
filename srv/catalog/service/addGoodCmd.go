package service

type AddGoodCmd struct {
	id          string
	name        string
	description string
}

func NewAddGoodCmd(id string, name string, description string) *AddGoodCmd {
	return &AddGoodCmd{id: id, name: name, description: description}
}

func (agc *AddGoodCmd) GetId() string {
	return agc.id
}

func (agc *AddGoodCmd) GetName() string {
	return agc.name
}

func (agc *AddGoodCmd) GetDescription() string {
	return agc.description
}
