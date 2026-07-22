package controller_admin_authentication

import (
	"github.com/gin-gonic/gin"

	inputApplicationHttp "example/internal/input/application/http"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
)

type AuthenticatorHandler struct {
	*inputApplicationHttp.AbstractHandler
	AuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase
}

// NewUserHandler 構造函數 (Go 的慣用法)，
// 相当 PHP 的 __construct()

func NewAuthenticatorHandler(oAbstractHandler *inputApplicationHttp.AbstractHandler, oAuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler:      oAbstractHandler,
		AuthenticatorUsecase: oAuthenticatorUsecase,
	}
}

func (oSelf *AuthenticatorHandler) SignIn(oContext *gin.Context) {

	sParamName := oContext.PostForm("param.name")
	sParamPassword := oContext.PostForm("param.password")

	sAuthorization, oErr := oSelf.AuthenticatorUsecase.SignIn(
		sParamName,
		sParamPassword,
	)

	if oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, oErr.Error(), struct{}{}, "")
		return
	}

	oSelf.Response.SetWithNext(oContext, 200, 1, "成功登入", struct{}{}, sAuthorization)

}
