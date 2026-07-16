package controller_admin_authentication

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"example/internal/bootstrap"
	inputHttpAdmin "example/internal/input/application/http/admin"
	outputPortModel "example/internal/output/port/any/model"
	"example/internal/utility"
)

type AuthenticatorHandler struct {
	*inputHttpAdmin.AbstractHandler
	AdminUserModelRepository outputPortModel.AdminUserRepository
}

// NewUserHandler 構造函數 (Go 的慣用法)，
// 相当 PHP 的 __construct()

func NewAuthenticatorHandler(oAbstractHandler *inputHttpAdmin.AbstractHandler, oAdminUserModelRepository outputPortModel.AdminUserRepository) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler:          oAbstractHandler,
		AdminUserModelRepository: oAdminUserModelRepository,
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

	oAdminUser, oErr := oSelf.AdminUserModelRepository.ShowOneByName(
		sParamName,
	)

	if oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			// 這是正確的判斷方式：處理找不到紀錄的情況
			oSelf.Response.Set(oContext, 200, -2, "AdminUser 不存在", struct{}{}, "")
			return
		}
	}

	if oAdminUser == nil {
		oSelf.Response.Set(oContext, 200, -2, "AdminUser 不存在", struct{}{}, "")
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
	}

	oSelf.Response.SetWithNext(oContext, 200, 1, "成功登入", struct{}{}, sAuthorization)

}
