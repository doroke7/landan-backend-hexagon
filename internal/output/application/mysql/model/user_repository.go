package mysql

import (
	"errors"

	"example/internal/output/port/any/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"example/internal/domain"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) port.UserRepository {
	return &UserRepository{db: db}
}

func (oSelf *UserRepository) AddOne(user *domain.User) error {
	return oSelf.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(user).Error
}

func (oSelf *UserRepository) ShowOneById(id int) (*domain.User, error) {
	var user domain.User
	if err := oSelf.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return &user, nil
}
