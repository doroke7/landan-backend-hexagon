package any

import (
	"errors"
	"testing"

	domain "example/internal/domain"
)

// fakeAppUserRepository 模擬真實 mysql repository 的行為（balance + amount），
// 不是單純回傳寫死的值，這樣測試斷言（20+10=30）才有意義、不會誤導人。
type fakeAppUserRepository struct {
	gotId        uint
	gotAmount    uint
	startBalance uint
	retErr       error
}

func (f *fakeAppUserRepository) IncreaseBalance(id uint, amount uint) (*domain.AppUser, error) {
	f.gotId = id
	f.gotAmount = amount
	if f.retErr != nil {
		return nil, f.retErr
	}
	return &domain.AppUser{Id: id, Balance: f.startBalance + amount}, nil
}

func TestAppUserUsecase_IncreaseBalance(t *testing.T) {
	t.Run("成功時把 amount 加到起始 balance 上，並把 id/amount 原封不動轉過去", func(t *testing.T) {
		oRepo := &fakeAppUserRepository{startBalance: 20}
		oUsecase := NewAppUserUsecase(oRepo)

		oUser, err := oUsecase.IncreaseBalance(1, 10)

		if err != nil {
			t.Fatalf("預期沒有錯誤，卻收到: %v", err)
		}
		if oRepo.gotId != 1 || oRepo.gotAmount != 10 {
			t.Fatalf("repository 收到的參數不對，got id=%d amount=%d", oRepo.gotId, oRepo.gotAmount)
		}
		if oUser.Balance != 30 {
			t.Fatalf("預期 20+10=30，got balance=%d", oUser.Balance)
		}
	})

	t.Run("repository 失敗時，usecase 要原樣把錯誤往上傳", func(t *testing.T) {
		oExpectedErr := errors.New("db down")
		oRepo := &fakeAppUserRepository{retErr: oExpectedErr}
		oUsecase := NewAppUserUsecase(oRepo)

		_, err := oUsecase.IncreaseBalance(1, 10)

		if !errors.Is(err, oExpectedErr) {
			t.Fatalf("預期錯誤被原樣傳回，got: %v", err)
		}
	})
}
