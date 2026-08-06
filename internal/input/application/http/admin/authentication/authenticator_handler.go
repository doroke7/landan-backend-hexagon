package controller_admin_authentication

import (
	"github.com/gin-gonic/gin"

	bootstrap "example/bootstrap"
	inputApplicationHttp "example/internal/input/application/http"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
	"example/pkg"
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
		bootstrap.CONFIG.SERVICES.HTTP.ADMIN.JWT.SECRET,
	)

	if oErr != nil {

		if oDefaultError, bOk := oErr.(*pkg.DefaultError); bOk {

			oSelf.Response.Set(oContext, 200, int(oDefaultError.Code), oDefaultError.Message, struct{}{}, "")
			return
		}

		oSelf.Response.Set(oContext, 200, -3, oErr.Error(), struct{}{}, "")
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "成功登入", struct{}{}, sAuthorization)

}
