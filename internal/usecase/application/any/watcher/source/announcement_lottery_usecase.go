package source

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyWatcherSource "example/internal/usecase/port/any/watcher/source"
	pkg "example/pkg"

	"go.uber.org/zap"
)

type AnnouncementLotteryUsecase struct {
	*AbstractUsecase
	outputPortAnyModel.LotteryRepository
}

func NewAnnouncementLotteryUsecase(oAbstractUsecase *AbstractUsecase, oLotteryRepository outputPortAnyModel.LotteryRepository) usecasePortAnyWatcherSource.AnnouncementLotteryUsecase {
	return &AnnouncementLotteryUsecase{
		AbstractUsecase:   oAbstractUsecase,
		LotteryRepository: oLotteryRepository,
	}
}

// Watch 收到一筆開獎資料，用 Round 當 key 落地存起來。
func (oSelf *AnnouncementLotteryUsecase) Watch(oLottery *domain.Lottery) error {

	if _, err := oSelf.LotteryRepository.EditOneByKey(oLottery, oLottery.Round); err != nil {
		pkg.Logger(pkg.Client).Error("儲存開獎資料失敗",
			zap.Uint("id", oLottery.Id),
			zap.String("round", oLottery.Round),
			zap.Error(err),
		)
		return err
	}

	pkg.Logger(pkg.Client).Info("收到開獎資料",
		zap.Uint("id", oLottery.Id),
		zap.String("round", oLottery.Round),
		zap.String("numbers", oLottery.Numbers),
	)

	return nil
}
