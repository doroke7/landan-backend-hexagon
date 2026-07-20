package output_port_any

import (
	domain "example/internal/domain"
)

type AdminUserRepository interface {
	ShowOneByName(name string) (*domain.AdminUser, error)
}
