package catalogCommon

type Good struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ID          string `json:"id"` //it is a converted uuid.UUID
}

func NewGood(ID string, name string, description string) *Good {
	return &Good{name, description, ID}
}

func (g Good) GetID() string {
	return g.ID
}

func (g Good) GetName() string {
	return g.Name
}

func (g Good) GetDescription() string {
	return g.Description
}

func (g *Good) SetDescription(newDescription string) error {
	if newDescription == "" {
		return CustomError{"Description is empty"}
	}
	g.Description = newDescription
	return nil
}

func (g *Good) SetName(newName string) error {
	if newName == "" {
		return CustomError{"Name is empty"}
	}
	g.Name = newName
	return nil
}
