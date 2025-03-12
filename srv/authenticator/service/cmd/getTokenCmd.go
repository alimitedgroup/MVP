package servicecmd

type GetTokenCmd struct {
	username string
}

func NewGetTokenCmd(username string) *GetTokenCmd {
	return &GetTokenCmd{username: username}
}

func (gtc *GetTokenCmd) GetUsername() string {
	return gtc.username
}
