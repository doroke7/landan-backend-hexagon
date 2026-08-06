package any

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
)

type AppUserUsecase struct {
	outputPortAnyModel.AppUserRepository
}

func NewAppUserUsecase(oAppUserRepository outputPortAnyModel.AppUserRepository) usecasePortAnyAdminResource.AppUserUsecase {
	return &AppUserUsecase{
		AppUserRepository: oAppUserRepository,
	}
}

func (oSelf *AppUserUsecase) IncreaseBalance(id uint, amount uint) (*domain.AppUser, error) {
	return oSelf.AppUserRepository.IncreaseBalance(id, amount)
}
