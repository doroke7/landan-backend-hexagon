package announcement

import (
	domain "example/internal/domain"
)

type LotteryUsecase interface {
	WatchOneByKey(key string) (*domain.Lottery, error)
}
