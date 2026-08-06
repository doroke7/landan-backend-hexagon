package cache

import (
	"fmt"

	"example/internal/domain"
	"example/internal/helper"
	"example/internal/output/mysql"
	"example/internal/output/port"
)

// UserRepository 是裝飾器（Decorator）：包住 mysql.UserRepository，
// 對外一樣實作 port.UserRepository，寫入時順便寫一份到 redis（write-through），
// 讀取時先查 redis，沒有才 fallback 查 mysql。usecase 完全不知道背後多包了一層快取。
type UserRepository struct {
	next *mysql.UserRepository // 5. 這樣看來 六角框架 的輸出 支持 多種ouput 嵌套（如 mysql 嵌套 redis 輸出）

	cacheHelper *helper.CacheHelper
}

func NewUserRepository(oInner *mysql.UserRepository, oCacheHelper *helper.CacheHelper) port.UserRepository {
	return &UserRepository{next: oInner, cacheHelper: oCacheHelper}
}

func (oSelf *UserRepository) AddOne(user *domain.User) error {
	if err := oSelf.next.AddOne(user); err != nil {
		return err
	}

	return oSelf.cacheHelper.EvictCache(oSelf.cacheKey(user.ID))
}

func (oSelf *UserRepository) ShowOneById(id int) (*domain.User, error) {
	var user domain.User
	if err := oSelf.cacheHelper.ReadCache(oSelf.cacheKey(id), &user); err == nil {
		return &user, nil
	}

	oUser, err := oSelf.next.ShowOneById(id)
	if err != nil {
		return nil, err
	}

	_ = oSelf.cacheHelper.WriteCache(oSelf.cacheKey(oUser.ID), oUser)
	return oUser, nil
}

func (oSelf *UserRepository) cacheKey(id int) string {
	return fmt.Sprintf("user:%d", id)
}
