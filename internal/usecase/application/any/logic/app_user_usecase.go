package usecase

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
)

type AppUserUsecase struct {
	*AbstractUsecase
	outputPortAnyLogic.AppUserRepository
}

func NewAppUserUsecase(oAppUserRepository outputPortAnyLogic.AppUserRepository, oAbstractUsecase *AbstractUsecase) usecasePortAnyLogic.AppUserUsecase {
	return &AppUserUsecase{
		AbstractUsecase:   oAbstractUsecase,
		AppUserRepository: oAppUserRepository,
	}
}

func (oSelf *AppUserUsecase) AddAppUser(oAppUser *domain.AppUser) (*domain.AppUser, error) {

	return oAppUser, nil
}
