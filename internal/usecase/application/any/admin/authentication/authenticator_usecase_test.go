package usecase

import (
	"errors"
	"testing"

	domain "example/internal/domain"
)

type fakeAdminUserRepository struct {
	gotName string
	retUser *domain.AdminUser
	retErr  error
}

func (f *fakeAdminUserRepository) ShowOneByName(name string) (*domain.AdminUser, error) {
	f.gotName = name
	return f.retUser, f.retErr
}

func TestAuthenticatorUsecase_ShowOneByName(t *testing.T) {
	t.Run("找到帳號時原封不動回傳，並把 name 正確轉過去", func(t *testing.T) {
		oExpected := &domain.AdminUser{Id: 1, Name: "tom", Password: "hashed"}
		oRepo := &fakeAdminUserRepository{retUser: oExpected}
		oUsecase := NewAuthenticatorUsecase(oRepo, &AbstractUsecase{})

		oAdminUser, err := oUsecase.ShowOneByName("tom")

		if err != nil {
			t.Fatalf("預期沒有錯誤，卻收到: %v", err)
		}
		if oRepo.gotName != "tom" {
			t.Fatalf("repository 收到的 name 不對，got %q", oRepo.gotName)
		}
		if oAdminUser != oExpected {
			t.Fatalf("預期原封不動回傳 repository 的結果，got %+v", oAdminUser)
		}
	})

	t.Run("repository 沒回傳錯誤、但也沒找到帳號（nil, nil）時，要轉成明確的 not found 錯誤", func(t *testing.T) {
		oRepo := &fakeAdminUserRepository{retUser: nil, retErr: nil}
		oUsecase := NewAuthenticatorUsecase(oRepo, &AbstractUsecase{})

		oAdminUser, err := oUsecase.ShowOneByName("nobody")

		if oAdminUser != nil {
			t.Fatalf("預期回傳 nil，got %+v", oAdminUser)
		}
		if err == nil {
			t.Fatal("預期回傳 not found 錯誤，卻沒有錯誤")
		}
	})

	t.Run("repository 失敗時，usecase 要原樣把錯誤往上傳", func(t *testing.T) {
		oExpectedErr := errors.New("db down")
		oRepo := &fakeAdminUserRepository{retErr: oExpectedErr}
		oUsecase := NewAuthenticatorUsecase(oRepo, &AbstractUsecase{})

		_, err := oUsecase.ShowOneByName("tom")

		if !errors.Is(err, oExpectedErr) {
			t.Fatalf("預期錯誤被原樣傳回，got: %v", err)
		}
	})
}
