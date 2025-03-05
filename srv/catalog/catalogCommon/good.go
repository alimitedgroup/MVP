package catalogCommon

type Good struct {
	name        string
	description string
	ID          string //it is a converted uuid.UUID
}

func NewGood(ID string, name string, description string) *Good {
	return &Good{name, description, ID}
}

func (g Good) GetID() string {
	return g.ID
}

func (g Good) GetName() string {
	return g.name
}

func (g Good) GetDescription() string {
	return g.description
}

func (g *Good) SetDescription(newDescription string) error {
	if newDescription == "" {
		return CustomError{"Description is empty"}
	}
	g.description = newDescription
	return nil
}

func (g *Good) SetName(newName string) error {
	if newName == "" {
		return CustomError{"Name is empty"}
	}
	g.name = newName
	return nil
}
