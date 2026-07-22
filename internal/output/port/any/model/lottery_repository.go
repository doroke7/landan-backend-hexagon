package port

import (
	domain "example/internal/domain"
)

type LotteryRepository interface {
	EditOneByKey(oLottery *domain.Lottery, sKey string) (*domain.Lottery, error)
	WatchOneByKey(sKey string) (*domain.Lottery, error)
}
