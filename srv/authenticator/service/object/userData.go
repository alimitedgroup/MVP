package serviceobject

type UserData struct {
	username string
}

func NewUserData(username string) *UserData {
	return &UserData{username: username}
}

func (ud *UserData) GetUsername() string {
	return ud.username
}
