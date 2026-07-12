package port

import "example/internal/domain"

type UserUsecase interface {
	AddUserByName(name string) (*domain.User, error)
	ShowUserById(id int) (*domain.User, error)
}
