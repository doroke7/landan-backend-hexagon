package cron

import (
	port "example/internal/usecase/port/cron"
)

type AppUserUsecase struct {
}

func NewAppUserUsecase() port.AppUserUsecase {
	return &AppUserUsecase{}
}

func (oSelf *AppUserUsecase) IncreaseBalance() {

}
