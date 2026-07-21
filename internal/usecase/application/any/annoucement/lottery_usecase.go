package usecase

import (
	domain "example/internal/domain"

	// outputPortAnyLogic "example/internal/output/port/any/logic"
	usecasePortAnyAnnoucement "example/internal/usecase/port/any/announcement"
)

type LotteryUsecase struct {
	*AbstractUsecase
}

func NewLotteryUsecase(oAbstractUsecase *AbstractUsecase) usecasePortAnyAnnoucement.LotteryUsecase {
	return &LotteryUsecase{
		AbstractUsecase: oAbstractUsecase,
	}
}

// 這邊有一個小技巧
// 1. 我們先關心 usecase 的情況 是否符合介面的情況，不要那麼早的注入 repository
// 2. 先 直接 return 一個假的 domain hard-code 資料

func (oSelf *LotteryUsecase) WatchOneByKey(sKey string) (*domain.Lottery, error) {

	return &domain.Lottery{
		Id:      1,
		Round:   "2026-001-001",
		Time:    111111111111111111,
		Numbers: "7,97,72,53",
	}, nil
}
