package cron

import (
	"example/internal/usecase/port/model"
	pkg "example/pkg"

	"go.uber.org/zap"
)

type createUserMessage struct {
	Name string `json:"name"`
}

type UserCron struct {
	*AbstractHandler
	userUsecase port.UserUsecase // UserUsecase 是 driving port，是每個 handler 各自要注入的業務依賴，不是「抽象共用的技術基礎設施」，不能塞進 AbstractHandler
}

func NewUserCron(usecase port.UserUsecase, oAbstractHandler *AbstractHandler) (*UserCron, error) {
	return &UserCron{AbstractHandler: oAbstractHandler, userUsecase: usecase}, nil
}

// AddUser 是 "user.create" queue 對應的處理方法，職責跟 http/grpc handler 的
// AddUser 一樣單純：只管收到的訊息怎麼轉呼叫 usecase，佇列宣告／消費迴圈都交給 register.ConsumerRouter。
func (oSelf *UserCron) AddUser() {
	var payload = &createUserMessage{
		Name: "Joe",
	}

	if _, err := oSelf.userUsecase.AddUserByName(payload.Name); err != nil {
		pkg.Logger(pkg.Cron).Error("AddUser 失敗",
			zap.String("name", payload.Name),
			zap.Error(err),
		)
		return
	}

	pkg.Logger(pkg.Cron).Info("AddUser 成功",
		zap.String("name", payload.Name),
	)
}
