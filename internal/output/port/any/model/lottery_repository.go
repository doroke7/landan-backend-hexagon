package port

import (
	domain "example/internal/domain"
)

type LotteryRepository interface {
	WatchOneByKey(sKey string) (*domain.Lottery, error)
}
