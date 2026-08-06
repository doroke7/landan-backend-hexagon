package port

import "example/domain"

type UserUsecase interface {
	CreateUser(name string) (*domain.User, error)
	GetUser(id int) (*domain.User, error)
}
