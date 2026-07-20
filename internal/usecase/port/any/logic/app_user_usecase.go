package any

import (
	"example/internal/domain"
)

type AppUserUsecase interface {
	AddAppUser(oAppUser *domain.AppUser) (*domain.AppUser, error)
}
