package authentication

import (
	"go.uber.org/zap"

	pkg "example/pkg"

	inputApplicationCron "example/internal/input/application/cron"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
)

type AuthenticatorHandler struct {
	*inputApplicationCron.AbstractHandler
	AuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase
}

func NewAuthenticatorHandler(oAuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase, oAbstractHandler *inputApplicationCron.AbstractHandler) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler:      oAbstractHandler,
		AuthenticatorUsecase: oAuthenticatorUsecase,
	}
}

func (oSelf *AuthenticatorHandler) SignIn() {
	sAuthorization, err := oSelf.AuthenticatorUsecase.SignIn("tom", "secret")
	if err != nil {
		pkg.Logger(pkg.Cron).Error("SignIn 失敗",
			zap.Error(err),
		)
		return
	}

	pkg.Logger(pkg.Cron).Info("SignIn 成功",
		zap.String("authorization", sAuthorization),
	)
}
