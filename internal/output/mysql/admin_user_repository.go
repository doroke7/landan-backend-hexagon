package mysql

import (
	"errors"

	"gorm.io/gorm"

	"example/internal/domain"
	"example/internal/output/port"
)

type AdminUserRepository struct {
	db *gorm.DB
}

func NewAdminUserRepository(db *gorm.DB) port.AdminUserRepository {
	return &AdminUserRepository{db: db}
}

func (oSelf *AdminUserRepository) ShowOneByName(sName string) (*domain.AdminUser, error) {
	var oAdminUser domain.AdminUser

	if err := oSelf.db.Where("name = ?", sName).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &oAdminUser, nil
}
