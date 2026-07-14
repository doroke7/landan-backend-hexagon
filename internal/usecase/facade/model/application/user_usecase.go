package usecase

import (
	"errors"

	"example/internal/domain"
	outputPort "example/internal/output/port"
	inputPort "example/internal/usecase/facade/model/port"
)

type UserUsecase struct {
	*AbstractUsecase
	outputPort.UserRepository
}

func NewUserUsecase(oUserRepository outputPort.UserRepository, oAbstractUsecase *AbstractUsecase) inputPort.UserUsecase {
	return &UserUsecase{
		AbstractUsecase: oAbstractUsecase,
		UserRepository:  oUserRepository,
	}
}

func (oSelf *UserUsecase) AddUserByName(name string) (*domain.User, error) {

	user := &domain.User{
		Name: name,
	}

	if err := oSelf.UserRepository.AddOne(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (oSelf *UserUsecase) ShowUserById(id int) (*domain.User, error) {

	user, err := oSelf.UserRepository.ShowOneById(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
