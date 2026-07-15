package usecase

import (
	"errors"

	"example/internal/domain"
	inputPort "example/internal/usecase/resource/model/port"
	outputPort "example/internal/output/port"
)

type AdminUserUsecase struct {
	*AbstractUsecase
	outputPort.AdminUserRepository
}

func NewAdminUserUsecase(oAminUserRepository outputPort.AdminUserRepository, oAbstractUsecase *AbstractUsecase) inputPort.AdminUserUsecase {
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
