package service

type AddChangeGoodCmd struct {
	id          string
	name        string
	description string
}

func NewAddChangeGoodCmd(id string, name string, description string) *AddChangeGoodCmd {
	return &AddChangeGoodCmd{id: id, name: name, description: description}
}

func (agc *AddChangeGoodCmd) GetId() string {
	return agc.id
}

func (agc *AddChangeGoodCmd) GetName() string {
	return agc.name
}

func (agc *AddChangeGoodCmd) GetDescription() string {
	return agc.description
}
