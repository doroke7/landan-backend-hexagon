package usecase

import (
	bootstrap "example/bootstrap"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
	utility "example/internal/utility"
	pkg "example/pkg"
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

func (oSelf *AuthenticatorUsecase) SignIn(sName string, sPassword string, sSecret string) (string, error) {

	if sName == "" {
		return "", pkg.NewDefaultError("name 不能為空", -1, 200)

	}

	if sPassword == "" {
		return "", pkg.NewDefaultError("password 不能為空", -1, 200)
	}

	oAdminUser, err := oSelf.AdminUserRepository.ShowOneByName(sName)

	if err != nil {
		return "", err
	}
	if oAdminUser == nil {
		return "", pkg.NewDefaultError(sName+" 不存在", -2, 200)
	}

	sMd5 := utility.Md5(sPassword + bootstrap.CONFIG.TABLE.ADMIN_USER.PASSWORD)
	if oAdminUser.Password != sMd5 {
		return "", pkg.NewDefaultError("密碼錯誤", -2, 200)
	}

	sAuthorization, err := oSelf.JwtHelper.Generate(int64(oAdminUser.Id), 0, map[string]any{}, sSecret)
	if err != nil {
		return "", pkg.NewDefaultError("JWT 產生失敗", -2, 200)
	}

	return sAuthorization, nil
}
