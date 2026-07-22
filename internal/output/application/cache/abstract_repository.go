package cache

import (
	helper "example/internal/helper"
)

// AbstractRepository 放 cache 這個 output adapter 共用的依賴——CacheHelper，
// 跟 mysql.AbstractRepository 持有 *gorm.DB 是同一種角色。
type AbstractRepository struct {
	CacheHelper *helper.CacheHelper
}

func NewAbstractRepository(oCacheHelper *helper.CacheHelper) *AbstractRepository {
	return &AbstractRepository{
		CacheHelper: oCacheHelper,
	}
}
