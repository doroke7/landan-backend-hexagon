package output_application_mysql

import (
	"errors"

	"gorm.io/gorm"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type AdminUserRepository struct {
	*AbstractRepository
}

func NewAdminUserRepository(oAbstractRepository *AbstractRepository) outputPortAnyModel.AdminUserRepository {
	return &AdminUserRepository{
		AbstractRepository: oAbstractRepository,
	}
}

func (oSelf *AdminUserRepository) ShowOneByName(sName string) (*domain.AdminUser, error) {
	var oAdminUser domain.AdminUser

	if err := oSelf.db.Where("name = ?", sName).First(&oAdminUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return &oAdminUser, nil
}
