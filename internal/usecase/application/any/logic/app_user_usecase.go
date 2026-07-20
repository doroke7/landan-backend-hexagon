package usecase

import (
	"example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	usecasePortResourceLogic "example/internal/usecase/port/any/logic"
)

type AppUserUsecase struct {
	*AbstractUsecase
	outputPortAnyLogic.AppUserRepository
}

func NewAppUserUsecase(oAppUserRepository outputPortAnyLogic.AppUserRepository, oAbstractUsecase *AbstractUsecase) usecasePortResourceLogic.AppUserUsecase {
	return &AppUserUsecase{
		AbstractUsecase:   oAbstractUsecase,
		AppUserRepository: oAppUserRepository,
	}
}

func (oSelf *AppUserUsecase) AddAppUser(oAdminUser *domain.AdminUser) (*domain.AdminUser, error) {

	return oAdminUser, nil
}
