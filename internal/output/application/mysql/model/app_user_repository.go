package output_application

import (
	"gorm.io/gorm"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type AppUserRepository struct {
	*AbstractRepository
}

func NewAppUserRepository(oAbstractRepository *AbstractRepository) outputPortAnyModel.AppUserRepository {
	return &AppUserRepository{
		AbstractRepository: oAbstractRepository,
	}
}

func (oSelf *AppUserRepository) IncreaseBalance(id uint, amount uint) (*domain.AppUser, error) {
	if err := oSelf.db.Model(&domain.AppUser{}).
		Where("id = ?", id).
		UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
		return nil, err
	}

	var oAppUser domain.AppUser
	if err := oSelf.db.Where("id = ?", id).First(&oAppUser).Error; err != nil {
		return nil, err
	}

	return &oAppUser, nil
}
