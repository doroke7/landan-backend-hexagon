package cron

import (
	port "example/internal/usecase/port/cron"
)

type AppUserHandler struct {
	*AbstractHandler
	appUserUsecase port.AppUserUsecase
}

func NewAppUserHandler(oAppUserUsecase port.AppUserUsecase, oAbstractHandler *AbstractHandler) *AppUserHandler {
	return &AppUserHandler{
		AbstractHandler: oAbstractHandler,
		appUserUsecase:  oAppUserUsecase,
	}
}

func (oSelf *AppUserHandler) IncreaseBalance() {
	oSelf.appUserUsecase.IncreaseBalance()
}
