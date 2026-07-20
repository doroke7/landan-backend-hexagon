package cron

import (
	inputCron "example/internal/input/application/cron"
	port "example/internal/usecase/port/cron/admin/resource"
	pkg "example/pkg"

	"go.uber.org/zap"
)

type AppUserHandler struct {
	*inputCron.AbstractHandler
	appUserUsecase port.AppUserUsecase
}

func NewAppUserHandler(oAppUserUsecase port.AppUserUsecase, oAbstractHandler *inputCron.AbstractHandler) *AppUserHandler {
	return &AppUserHandler{
		AbstractHandler: oAbstractHandler,
		appUserUsecase:  oAppUserUsecase,
	}
}

func (oSelf *AppUserHandler) IncreaseBalance() {
	oAppUser, err := oSelf.appUserUsecase.IncreaseBalance()
	if err != nil {
		pkg.Logger(pkg.Cron).Error("IncreaseBalance 失敗",
			zap.Error(err),
		)
		return
	}

	pkg.Logger(pkg.Cron).Info("IncreaseBalance 成功",
		zap.Uint("id", oAppUser.Id),
		zap.Uint("balance", oAppUser.Balance),
	)
}
