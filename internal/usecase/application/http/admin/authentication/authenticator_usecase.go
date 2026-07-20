package usecase

import (
	"errors"

	"example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortReouseceModel "example/internal/usecase/port/resource/model"
)

type AuthenticatorUsecase struct {
	*AbstractUsecase
	outputPortAnyModel.AdminUserRepository
}

func NewAuthenticatorUsecase(oAminUserRepository outputPortAnyModel.AdminUserRepository, oAbstractUsecase *AbstractUsecase) usecasePortReouseceModel.AdminUserUsecase {
	return &AuthenticatorUsecase{
		AbstractUsecase:     oAbstractUsecase,
		AdminUserRepository: oAminUserRepository,
	}
}

func (oSelf *AuthenticatorUsecase) ShowOneByName(sName string) (*domain.AdminUser, error) {

	oAdminUser, err := oSelf.AdminUserRepository.ShowOneByName(sName)
	if err != nil {
		return nil, err
	}

	if oAdminUser == nil {
		return nil, errors.New("AdminUser not found")
	}

	return oAdminUser, nil
}
