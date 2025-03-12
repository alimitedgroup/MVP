package serviceobject

type UserData struct {
	username string
	role     string
}

func NewUserData(username string, role string) *UserData {
	return &UserData{username: username, role: role}
}

func (ud *UserData) GetUsername() string {
	return ud.username
}

func (ud *UserData) GetRole() string {
	return ud.role
}
