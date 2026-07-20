package output_port

import (
	domain "example/internal/domain"
)

type AppUserRepository interface {
	AddAppUser(oAppUser *domain.AppUser) (*domain.AppUser, error)
}
