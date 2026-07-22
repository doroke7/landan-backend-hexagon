package source

import (
	domain "example/internal/domain"
	usecasePortAnyWatcherSource "example/internal/usecase/port/any/watcher/source"
	pkg "example/pkg"

	"go.uber.org/zap"
)

type AnnouncementLotteryUsecase struct {
	*AbstractUsecase
}

func NewAnnouncementLotteryUsecase(oAbstractUsecase *AbstractUsecase) usecasePortAnyWatcherSource.AnnouncementLotteryUsecase {
	return &AnnouncementLotteryUsecase{
		AbstractUsecase: oAbstractUsecase,
	}
}

// HandleDraw 收到一筆開獎資料要做什麼，業務邏輯寫在這裡——先只記 log，
// 之後如果要落地存檔，再幫這裡注入一個真正的 output repository。
func (oSelf *AnnouncementLotteryUsecase) Watch(oLottery *domain.Lottery) error {

	pkg.Logger(pkg.Client).Info("收到開獎資料",
		zap.Uint("id", oLottery.Id),
		zap.String("round", oLottery.Round),
		zap.String("numbers", oLottery.Numbers),
	)

	return nil
}
