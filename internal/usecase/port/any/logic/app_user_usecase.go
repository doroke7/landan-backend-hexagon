package any

import (
	"example/internal/domain"
)

type AppUserUsecase interface {
	AddAppUser(oAdminUser *domain.AdminUser) (*domain.AdminUser, error)
}
