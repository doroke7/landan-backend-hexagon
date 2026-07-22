package model

import (
	domain "example/internal/domain"
	outputApplicationCache "example/internal/output/application/cache"
	outputPortAnyModel "example/internal/output/port/any/model"
)

// LotteryRepository 直接讀寫 redis，不包其他 repository。
type LotteryRepository struct {
	*outputApplicationCache.AbstractRepository
}

func NewLotteryRepository(oAbstractRepository *outputApplicationCache.AbstractRepository) outputPortAnyModel.LotteryRepository {
	return &LotteryRepository{
		AbstractRepository: oAbstractRepository,
	}
}

// WatchOneByKey 是讀：直接讀 redis，沒有就回傳錯誤（不再自己生資料）。
func (oSelf *LotteryRepository) WatchOneByKey(sKey string) (*domain.Lottery, error) {
	var oLottery domain.Lottery
	if err := oSelf.CacheHelper.ReadCache(oSelf.cacheKey(sKey), &oLottery); err != nil {
		return nil, err
	}

	return &oLottery, nil
}

// EditOneByKey 是寫：把呼叫端給的 oLottery 寫進 redis。
func (oSelf *LotteryRepository) EditOneByKey(oLottery *domain.Lottery, sKey string) (*domain.Lottery, error) {
	if err := oSelf.CacheHelper.WriteCache(oSelf.cacheKey(sKey), oLottery); err != nil {
		return nil, err
	}

	return oLottery, nil
}

func (oSelf *LotteryRepository) cacheKey(sKey string) string {
	return "lottery:" + sKey
}
