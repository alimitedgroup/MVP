package out

type AuthenticationAdapter struct {
}

func (*AuthenticationAdapter) GetToken(username string) (UserToken, error) {

}
func (*AuthenticationAdapter) GetRole(token UserToken) (UserRole, error) {

}
