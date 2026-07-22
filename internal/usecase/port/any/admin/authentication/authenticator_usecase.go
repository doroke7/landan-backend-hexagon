package any

type AuthenticatorUsecase interface {
	SignIn(name string, password string) (authorization string, err error)
}
