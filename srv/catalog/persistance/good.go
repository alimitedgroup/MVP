package goodRepository

type good struct {
	name        string
	description string
	ID          string //it is a converted uuid.UUID
}

func NewGood(ID string, name string, description string) *good {
	return &good{name, description, ID}
}

func (g good) GetID() string {
	return g.ID
}

func (g good) GetName() string {
	return g.name
}

func (g good) GetDescription() string {
	return g.description
}

func (g *good) SetDescription(newDescription string) error {
	if newDescription == "" {
		return CustomError{"Description is empty"}
	}
	g.description = newDescription
	return nil
}

func (g *good) SetName(newName string) error {
	if newName == "" {
		return CustomError{"Name is empty"}
	}
	g.name = newName
	return nil
}
