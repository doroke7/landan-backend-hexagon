package output_port

import (
	domain "example/internal/domain"
)

type AdminUserRepository interface {
	ShowOneByName(name string) (*domain.AdminUser, error)
}
