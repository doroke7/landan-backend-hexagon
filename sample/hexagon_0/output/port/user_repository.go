package port

import "example/domain"

type UserRepository interface {
	AddOne(user *domain.User) error
	ShowOneById(id int) (*domain.User, error)
}
