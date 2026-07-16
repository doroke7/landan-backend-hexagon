package memory

import (
	"errors"

	"example/internal/domain"
	"example/internal/output/port/any/model"
)

type UserRepository struct {
	data map[int]*domain.User
}

func NewUserRepository() port.UserRepository {
	return &UserRepository{
		data: make(map[int]*domain.User),
	}
}

func (oSelf *UserRepository) AddOne(user *domain.User) error {
	oSelf.data[user.ID] = user
	return nil
}

func (oSelf *UserRepository) ShowOneById(id int) (*domain.User, error) {
	user, ok := oSelf.data[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return user, nil
}
