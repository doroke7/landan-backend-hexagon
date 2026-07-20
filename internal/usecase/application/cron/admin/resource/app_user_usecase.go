package cron

import (
	"example/internal/domain"
	outputPort "example/internal/output/port/any/model"
	port "example/internal/usecase/port/cron"
)

type AppUserUsecase struct {
	outputPort.AppUserRepository
}

func NewAppUserUsecase(oAppUserRepository outputPort.AppUserRepository) port.AppUserUsecase {
	return &AppUserUsecase{
		AppUserRepository: oAppUserRepository,
	}
}

func (oSelf *AppUserUsecase) IncreaseBalance() (*domain.AppUser, error) {
	return oSelf.AppUserRepository.IncreaseBalance(1, 10)
}
