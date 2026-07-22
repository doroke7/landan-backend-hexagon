package usecase

import (
	"errors"

	bootstrap "example/bootstrap"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
	utility "example/internal/utility"
)

type AuthenticatorUsecase struct {
	*AbstractUsecase
	outputPortAnyModel.AdminUserRepository
}

func NewAuthenticatorUsecase(oAminUserRepository outputPortAnyModel.AdminUserRepository, oAbstractUsecase *AbstractUsecase) usecasePortAnyAdminAuthentication.AuthenticatorUsecase {
	return &AuthenticatorUsecase{
		AbstractUsecase:     oAbstractUsecase,
		AdminUserRepository: oAminUserRepository,
	}
}

func (oSelf *AuthenticatorUsecase) SignIn(sName string, sPassword string) (string, error) {

	if sName == "" {
		return "", errors.New("name 不能為空")
	}

	if sPassword == "" {
		return "", errors.New("password 不能為空")
	}

	oAdminUser, err := oSelf.AdminUserRepository.ShowOneByName(sName)
	if err != nil {
		return "", err
	}

	if oAdminUser == nil {
		return "", errors.New("AdminUser not found")
	}

	sMd5 := utility.Md5(sPassword + bootstrap.CONFIG.TABLE.ADMIN_USER.PASSWORD)
	if oAdminUser.Password != sMd5 {
		return "", errors.New("密碼錯誤")
	}

	sAuthorization, err := oSelf.JwtHelper.Generate(int64(oAdminUser.Id), 0, map[string]any{})
	if err != nil {
		return "", errors.New("JWT 產生失敗")
	}

	return sAuthorization, nil
}
