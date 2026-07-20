package any

import (
	"example/internal/domain"
	outputPort "example/internal/output/port/any/model"
	port "example/internal/usecase/port/any/admin/resource"
)

type AppUserUsecase struct {
	outputPort.AppUserRepository
}

func NewAppUserUsecase(oAppUserRepository outputPort.AppUserRepository) port.AppUserUsecase {
	return &AppUserUsecase{
		AppUserRepository: oAppUserRepository,
	}
}

func (oSelf *AppUserUsecase) IncreaseBalance(id uint, amount uint) (*domain.AppUser, error) {
	return oSelf.AppUserRepository.IncreaseBalance(id, amount)
}
