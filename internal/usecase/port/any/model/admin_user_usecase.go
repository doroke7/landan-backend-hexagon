package any

import (
	"example/internal/domain"
)

type AdminUserUsecase interface {
	ShowOneByName(name string) (*domain.AdminUser, error)
}
