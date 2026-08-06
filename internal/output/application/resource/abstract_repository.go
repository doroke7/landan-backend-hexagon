package mysql

import (
	"context"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractRepository 做法一致。
type AbstractRepository struct {
	Context context.Context
}

func NewAbstractRepository(oContext context.Context) *AbstractRepository {
	return &AbstractRepository{
		Context: oContext,
	}
}
