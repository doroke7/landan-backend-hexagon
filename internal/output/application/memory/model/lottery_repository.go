package model

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"
)

type LotteryRepository struct {
	*AbstractRepository
	startTime time.Time
}

func NewLotteryRepository(oAbstractRepository *AbstractRepository) outputPortAnyModel.LotteryRepository {
	return &LotteryRepository{
		AbstractRepository: oAbstractRepository,
		startTime:          time.Now(),
	}
}

func (oSelf *LotteryRepository) WatchOneByKey(sKey string) (*domain.Lottery, error) {

	iId := uint(time.Since(oSelf.startTime).Minutes()) + 1

	iCount := 4 // 產生幾個數字

	aNumbers := make([]string, iCount)

	for i := 0; i < iCount; i++ {
		n := rand.IntN(99) + 1 // 1~99
		aNumbers[i] = strconv.Itoa(n)
	}

	sNumbers := strings.Join(aNumbers, ",")

	return &domain.Lottery{
		Id:      iId,
		Round:   fmt.Sprintf("2026-%04d", iId),
		Time:    time.Now().UnixNano(),
		Numbers: sNumbers,
	}, nil

}
