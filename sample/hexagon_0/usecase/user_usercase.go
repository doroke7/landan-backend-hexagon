package usecase

import (
	"errors"

	"example/domain"
	inputPort "example/input/port"
	outputPort "example/output/port"
)

type userUsecase struct {
	UserRepository outputPort.UserRepository
}

func NewUserUsecase(oUserRepository outputPort.UserRepository) inputPort.UserUsecase {
	return &userUsecase{
		UserRepository: oUserRepository,
	}
}

func (oSelf *userUsecase) CreateUser(name string) (*domain.User, error) {

	user := &domain.User{
		Name: name,
	}

	if err := oSelf.UserRepository.AddOne(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (oSelf *userUsecase) GetUser(id int) (*domain.User, error) {

	user, err := oSelf.UserRepository.ShowOneById(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
