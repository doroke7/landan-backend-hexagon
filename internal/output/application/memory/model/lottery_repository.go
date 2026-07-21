package model

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type LotteryRepository struct {
	*AbstractRepository
}

func NewLotteryRepository(oAbstractRepository *AbstractRepository) outputPortAnyModel.LotteryRepository {
	return &LotteryRepository{
		AbstractRepository: oAbstractRepository,
	}
}

func (oSelf *LotteryRepository) WatchOneByKey(sKey string) (*domain.Lottery, error) {

	return &domain.Lottery{
		Id:      1,
		Round:   "2026-001-001",
		Time:    111111111111111111,
		Numbers: "7,97,72,53",
	}, nil

}
