package any

import "example/internal/domain"

type AppUserUsecase interface {
	IncreaseBalance(id uint, amount uint) (*domain.AppUser, error)
}
