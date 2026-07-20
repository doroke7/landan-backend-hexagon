package usecase

import (
	"errors"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
)

type AuthenticatorUsecase struct {
	*AbstractUsecase
	outputPortAnyModel.AdminUserRepository
}

func NewAuthenticatorUsecase(oAminUserRepository outputPortAnyModel.AdminUserRepository, oAbstractUsecase *AbstractUsecase) usecasePortAnyAdminAuthentication.AuthenticatorUsecase {
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
