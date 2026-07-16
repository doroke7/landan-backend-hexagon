package controller_admin_authentication

import (
	"github.com/gin-gonic/gin"

	HttpAdmin "example/internal/input/http/admin"
)

type AuthenticatorHandler struct {
	*HttpAdmin.AbstractHandler
}

// NewUserHandler 構造函數 (Go 的慣用法)，
// 相当 PHP 的 __construct()

func NewAuthenticatorHandler(oAbstractHandler *HttpAdmin.AbstractHandler) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler: oAbstractHandler,
	}
}

func (oSelf *AuthenticatorHandler) SignIn(oContext *gin.Context) {
	sParamName := oContext.PostForm("param.name")
	sParamPassword := oContext.PostForm("param.password")

	if sParamName == "" {
		oSelf.Response.SetWithNext(oContext, 200, -1, "name 不能為空", struct{}{}, "")
	}

	if sParamPassword == "" {
		oSelf.Response.SetWithNext(oContext, 200, -1, "password 不能為空", struct{}{}, "")
	}

	// oAdminUserModel, err := oSelf.ResourceClient.Model.AdminUser.ShowOneByName(
	// 	oContext,
	// 	&pbResourceModel.OneAdminUserRequest{
	// 		Name: sParamName,
	// 	},
	// )

	// if err != nil {
	// 	oSelf.Response.SetWithNext(oContext, 200, -2, "AdminUser 不能建立", struct{}{}, "")
	// }

	// oAdminUser, oErr := pkg.Cacheable[*modelDatabase.AdminUserModel]("", time.Hour, func() (*modelDatabase.AdminUserModel, error) {
	// 	return oAdminUserModel.ShowOneByName(sParamName)
	// }, sParamName)

	// if oErr != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		// 這是正確的判斷方式：處理找不到紀錄的情況
	// 		oSelf.Response.SetWithNext(oContext, 200, -2, "AdminUser 不存在", struct{}{}, "")
	// 	}
	// }

	// if oAdminUser == nil {
	// 	oSelf.Response.SetWithNext(oContext, 200, -2, "AdminUser 不存在", struct{}{}, "")
	// }

	// sMd5 := utility.Md5(sParamPassword + bootstrap.CONFIG.TABLE.ADMIN_USER.PASSWORD)

	// if oAdminUser.Password != sMd5 {
	// 	oSelf.Response.SetWithNext(oContext, 200, -2, "密碼錯誤", struct{}{}, "")
	// }

	// sAuthorization, oErrJwt := oSelf.JwtHelper.Generate(oAdminUser.Id, 0, map[string]any{})

	// if oErrJwt != nil {
	// 	oSelf.Response.SetWithNext(oContext, 200, -2, "JWT 產生失敗", struct{}{}, "")
	// }

	// oSelf.Response.SetWithNext(oContext, 200, 1, "成功登入", struct{}{}, sAuthorization)

}
