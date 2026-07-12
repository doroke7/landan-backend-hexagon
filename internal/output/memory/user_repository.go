package memory

import (
	"errors"

	"example/internal/domain"
	"example/internal/output/port"
)

type UserRepository struct {
	data map[int]*domain.User
}

func NewUserRepository() port.UserRepository {
	return &UserRepository{
		data: make(map[int]*domain.User),
	}
}

func (r *UserRepository) AddOne(user *domain.User) error {
	r.data[user.ID] = user
	return nil
}

func (r *UserRepository) ShowOneById(id int) (*domain.User, error) {
	user, ok := r.data[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return user, nil
}
