package mysql

import (
	"errors"

	"gorm.io/gorm"

	domain "example/internal/domain"
	outputApplicationMysql "example/internal/output/application/mysql"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type AdminUserRepository struct {
	*outputApplicationMysql.AbstractRepository
}

func NewAdminUserRepository(oAbstractRepository *outputApplicationMysql.AbstractRepository) outputPortAnyModel.AdminUserRepository {
	return &AdminUserRepository{
		AbstractRepository: oAbstractRepository,
	}
}

func (oSelf *AdminUserRepository) ShowOneByName(sName string) (*domain.AdminUser, error) {
	var oAdminUser domain.AdminUser

	if err := oSelf.DB.Where("name = ?", sName).First(&oAdminUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("沒有 " + sName + " 用戶")
		}
		return nil, err
	}

	return &oAdminUser, nil
}
