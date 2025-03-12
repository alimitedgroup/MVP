package servicecmd

type GetTokenCmd struct {
	username string
	role     string
}

func NewGetTokenCmd(username string, role string) *GetTokenCmd {
	return &GetTokenCmd{username: username, role: role}
}

func (gtc *GetTokenCmd) GetUsername() string {
	return gtc.username
}

func (gtc *GetTokenCmd) GetRole() string {
	return gtc.role
}
