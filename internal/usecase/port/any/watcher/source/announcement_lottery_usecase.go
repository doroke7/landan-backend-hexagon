package source

import (
	domain "example/internal/domain"
)

type AnnouncementLotteryUsecase interface {
	Watch(oLottery *domain.Lottery) error
}
