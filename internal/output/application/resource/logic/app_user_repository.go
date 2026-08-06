package resource

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type AppUserRepository struct {
}

func NewAppUserRepository() outputPortAnyModel.AppUserRepository {
	return &AppUserRepository{}
}

func (oSelf *AppUserRepository) AddAppUser(oAppUser *domain.AppUser) (*domain.AppUser, error) {

	return &domain.AppUser{
		Id:       1,
		Name:     "11",
		Password: "222222",
	}, nil
}

func (oSelf *AppUserRepository) IncreaseBalance(id uint, amount uint) (*domain.AppUser, error) {

	return &domain.AppUser{
		Id:      id,
		Balance: amount,
	}, nil
}
