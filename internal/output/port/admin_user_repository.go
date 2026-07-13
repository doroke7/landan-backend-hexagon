package port

import (
	"example/internal/domain"
)

type AdminUserRepository interface {
	ShowOneByName(name string) (*domain.AdminUser, error)
}
