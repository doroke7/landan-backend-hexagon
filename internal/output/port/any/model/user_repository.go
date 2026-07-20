package output_port_any

import (
	domain "example/internal/domain"
)

type UserRepository interface {
	AddOne(user *domain.User) error
	ShowOneById(id int) (*domain.User, error)
}
