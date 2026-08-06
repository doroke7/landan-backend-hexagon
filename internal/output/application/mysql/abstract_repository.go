package mysql

import (
	"context"

	pkg "example/pkg"

	"gorm.io/gorm"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractRepository 做法一致。
type AbstractRepository struct {
	DB      *gorm.DB
	Context context.Context
	*pkg.Aop
}

func NewAbstractRepository(oContext context.Context, oDb *gorm.DB, oAop *pkg.Aop) *AbstractRepository {

	return &AbstractRepository{
		DB:      oDb,
		Context: oContext,
		Aop:     oAop,
	}
}
