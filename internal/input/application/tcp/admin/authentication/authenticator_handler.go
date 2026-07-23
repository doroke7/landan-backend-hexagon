package authentication

import (
	"fmt"
	"strings"

	inputApplicationTcp "example/internal/input/application/tcp"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
)

type AuthenticatorHandler struct {
	*inputApplicationTcp.AbstractHandler
	AuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase
}

func NewAuthenticatorHandler(oAuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase, oAbstractHandler *inputApplicationTcp.AbstractHandler) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler:      oAbstractHandler,
		AuthenticatorUsecase: oAuthenticatorUsecase,
	}
}

// 不需要額外的轉接層。message body 格式固定「name:password」。
func (oSelf *AuthenticatorHandler) SignIn(aMessage []byte) ([]byte, error) {
	aParts := strings.SplitN(string(aMessage), ":", 2)
	if len(aParts) != 2 {
		return nil, fmt.Errorf("tcp: invalid SignIn message, expect \"name:password\"")
	}

	sAuthorization, err := oSelf.AuthenticatorUsecase.SignIn(aParts[0], aParts[1])
	if err != nil {
		return []byte(err.Error()), nil
	}

	return []byte(sAuthorization), nil
}
