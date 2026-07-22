package authentication

import (
	inputApplicationFacade "example/internal/input/application/facade"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
)

type AuthenticatorHandler struct {
	*inputApplicationFacade.AbstractHandler
	AuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase
}

func NewAuthenticatorHandler(oAuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase, oAbstractHandler *inputApplicationFacade.AbstractHandler) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler:      oAbstractHandler,
		AuthenticatorUsecase: oAuthenticatorUsecase,
	}
}

// SignIn 只負責轉呼叫 usecase，不知道自己是被 CLI 呼叫的——
// cobra.Command 的組裝、container 的建立時機，都交給 cmd/register 那層決定。
func (oSelf *AuthenticatorHandler) SignIn(sName string, sPassword string) (string, error) {
	return oSelf.AuthenticatorUsecase.SignIn(sName, sPassword)
}
