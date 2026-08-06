package usecase

import (
	domain "example/internal/domain"

	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyAnnoucement "example/internal/usecase/port/any/announcement"
)

type LotteryUsecase struct {
	*AbstractUsecase
	outputPortAnyModel.LotteryRepository
}

func NewLotteryUsecase(oAbstractUsecase *AbstractUsecase, oLotteryRepository outputPortAnyModel.LotteryRepository) usecasePortAnyAnnoucement.LotteryUsecase {
	return &LotteryUsecase{
		AbstractUsecase:   oAbstractUsecase,
		LotteryRepository: oLotteryRepository,
	}
}

// 這邊有一個小技巧
// 1. 我們先關心 usecase 的情況 是否符合介面的情況，不要那麼早的注入 repository
// 2. 先 直接 return 一個假的 domain hard-code 資料

func (oSelf *LotteryUsecase) WatchOneByKey(sKey string) (*domain.Lottery, error) {

	return oSelf.LotteryRepository.WatchOneByKey(sKey)

}
