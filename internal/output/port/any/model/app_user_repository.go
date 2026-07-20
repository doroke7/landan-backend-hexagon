package output_port

import (
	domain "example/internal/domain"
)

type AppUserRepository interface {
	IncreaseBalance(id uint, amount uint) (*domain.AppUser, error)
}
