package port

import (
	"example/internal/domain"
)

type AuthenticatorUsecase interface {
	ShowOneByName(name string) (*domain.AdminUser, error)
}
