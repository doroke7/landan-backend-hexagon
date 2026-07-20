package controller_admin_authentication

import (
	"github.com/gin-gonic/gin"

	bootstrap "example/bootstrap"
	utility "example/internal/utility"

	inputHttpAdmin "example/internal/input/application/http/admin"
	port "example/internal/usecase/port/http/admin/authentication"
)

type AuthenticatorHandler struct {
	*inputHttpAdmin.AbstractHandler
	AuthenticatorUsecase port.AuthenticatorUsecase
}

// NewUserHandler 構造函數 (Go 的慣用法)，
// 相当 PHP 的 __construct()

func NewAuthenticatorHandler(oAbstractHandler *inputHttpAdmin.AbstractHandler, oAuthenticatorUsecase port.AuthenticatorUsecase) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler:      oAbstractHandler,
		AuthenticatorUsecase: oAuthenticatorUsecase,
	}
}

func (oSelf *AuthenticatorHandler) SignIn(oContext *gin.Context) {
	sParamName := oContext.PostForm("param.name")
	sParamPassword := oContext.PostForm("param.password")

	if sParamName == "" {
		oSelf.Response.Set(oContext, 200, -1, "name 不能為空", struct{}{}, "")
		return
	}

	if sParamPassword == "" {
		oSelf.Response.Set(oContext, 200, -1, "password 不能為空", struct{}{}, "")
		return
	}

	oAdminUser, oErr := oSelf.AuthenticatorUsecase.ShowOneByName(
		sParamName,
	)

	if oErr != nil {
		oSelf.Response.Set(oContext, 200, -2, oErr.Error(), struct{}{}, "")
		return
	}

	sMd5 := utility.Md5(sParamPassword + bootstrap.CONFIG.TABLE.ADMIN_USER.PASSWORD)

	if oAdminUser.Password != sMd5 {
		oSelf.Response.SetWithNext(oContext, 200, -2, "密碼錯誤", struct{}{}, "")
		return
	}

	sAuthorization, oErr := oSelf.JwtHelper.Generate(int64(oAdminUser.Id), 0, map[string]any{})

	if oErr != nil {
		oSelf.Response.SetWithNext(oContext, 200, -2, "JWT 產生失敗", struct{}{}, "")
		return
	}

	oSelf.Response.SetWithNext(oContext, 200, 1, "成功登入", struct{}{}, sAuthorization)

}
