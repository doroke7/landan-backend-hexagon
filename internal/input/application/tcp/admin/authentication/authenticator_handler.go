package authentication

import (
	"strings"

	inputApplicationTcp "example/internal/input/application/tcp"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
	pkg "example/pkg"
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

// SignIn 簽名就是 pkg.TcpHandlerFunc，可以直接被 Tcp.HandleFunc 註冊。
// param 格式固定「name:password」。
func (oSelf *AuthenticatorHandler) SignIn(oReq pkg.TcpRequest) pkg.TcpResponse {
	aParts := strings.SplitN(oReq.Param, ":", 2)
	if len(aParts) != 2 {
		return pkg.TcpResponse{Code: -1, Message: "invalid param, expect \"name:password\""}
	}

	sAuthorization, err := oSelf.AuthenticatorUsecase.SignIn(aParts[0], aParts[1])
	if err != nil {
		return pkg.TcpResponse{Code: -1, Message: err.Error()}
	}

	return pkg.TcpResponse{Code: 1, Message: "成功登入", Result: sAuthorization}
}
