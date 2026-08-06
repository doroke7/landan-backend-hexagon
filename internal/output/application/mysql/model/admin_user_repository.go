package mysql

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	bootstrap "example/bootstrap"
	domain "example/internal/domain"
	outputApplicationMysql "example/internal/output/application/mysql"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type AdminUserRepository struct {
	*outputApplicationMysql.AbstractRepository
}

func NewAdminUserRepository(oAbstractRepository *outputApplicationMysql.AbstractRepository) outputPortAnyModel.AdminUserRepository {
	return &AdminUserRepository{
		AbstractRepository: oAbstractRepository,
	}
}

func (oSelf *AdminUserRepository) ShowOneByName(sName string) (*domain.AdminUser, error) {
	oCurrentContext, cCancel := context.WithTimeout(
		oSelf.Context,
		time.Duration(bootstrap.CONFIG.DATABASE.TIMEOUT)*time.Millisecond,
	)
	defer cCancel()

	var oAdminUser domain.AdminUser

	if err := oSelf.DB.WithContext(oCurrentContext).Where("name = ?", sName).First(&oAdminUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("資料不存在")
		}
		return nil, err
	}

	return &oAdminUser, nil
}

func (oSelf *AdminUserRepository) ShowOneById(iId uint) (*domain.AdminUser, error) {
	var oAdminUser domain.AdminUser
	sKey := oSelf.Aop.Key("AdminUser.SObI", iId)
	iTtl := oSelf.Aop.Ttl(30 * time.Minute)

	err := oSelf.Aop.Cacheable(sKey, iTtl, &oAdminUser, func() (interface{}, error) {
		oThisContext, cCancel := context.WithTimeout(
			oSelf.Context,
			time.Duration(bootstrap.CONFIG.DATABASE.TIMEOUT)*time.Millisecond,
		)
		defer cCancel()

		var oDbAdminUser domain.AdminUser
		if err := oSelf.DB.WithContext(oThisContext).First(&oDbAdminUser, iId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("資料不存在")
			}
			return nil, err
		}
		return oDbAdminUser, nil
	})

	return &oAdminUser, err
}
