package model

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"
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

	iCount := 4 // 產生幾個數字

	aNumbers := make([]string, iCount)

	for i := 0; i < iCount; i++ {
		n := rand.IntN(99) + 1 // 1~99
		aNumbers[i] = strconv.Itoa(n)
	}

	sNumbers := strings.Join(aNumbers, ",")

	return &domain.Lottery{
		Id:      1,
		Round:   "2026-001-001",
		Time:    time.Now().UnixNano(),
		Numbers: sNumbers,
	}, nil

}
