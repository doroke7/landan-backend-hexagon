package usecase

import (
	"errors"

	domain "example/internal/domain"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type AdminUserUsecase struct {
	*AbstractUsecase
	outputPortAnyModel.AdminUserRepository
}

func NewAdminUserUsecase(oAminUserRepository outputPortAnyModel.AdminUserRepository, oAbstractUsecase *AbstractUsecase) usecasePortAnyModel.AdminUserUsecase {
	return &AdminUserUsecase{
		AbstractUsecase:     oAbstractUsecase,
		AdminUserRepository: oAminUserRepository,
	}
}

func (oSelf *AdminUserUsecase) ShowOneByName(sName string) (*domain.AdminUser, error) {

	oAdminUser, err := oSelf.AdminUserRepository.ShowOneByName(sName)
	if err != nil {
		return nil, err
	}

	if oAdminUser == nil {
		return nil, errors.New("AdminUser not found")
	}

	return oAdminUser, nil
}
