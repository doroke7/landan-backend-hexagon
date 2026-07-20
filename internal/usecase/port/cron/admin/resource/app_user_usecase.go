package cron

import (
	"example/internal/domain"
)

type AppUserUsecase interface {
	IncreaseBalance() (*domain.AppUser, error)
}
