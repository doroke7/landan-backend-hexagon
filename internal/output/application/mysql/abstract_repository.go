package mysql

import (
	pkg "example/pkg"

	"gorm.io/gorm"
)

type AbstractRepository struct {
	DB *gorm.DB
	*pkg.Aop
}

func NewAbstractRepository(oDb *gorm.DB, oAop *pkg.Aop) *AbstractRepository {
	return &AbstractRepository{
		DB:  oDb,
		Aop: oAop,
	}
}
