package consumer

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"

	"example/internal/usecase/resource/port/model"
	pkg "example/pkg"

	"go.uber.org/zap"
)

type createUserMessage struct {
	Name string `json:"name"`
}

func NewUserConsumer(oUserUsecase port.UserUsecase, oAbstractHandler *AbstractHandler) (*UserConsumer, error) {
	return &UserConsumer{
		AbstractHandler: oAbstractHandler,
		userUsecase:     oUserUsecase,
	}, nil
}

type UserConsumer struct {
	*AbstractHandler
	userUsecase port.UserUsecase // UserUsecase 是 driving port，是每個 handler 各自要注入的業務依賴，不是「抽象共用的技術基礎設施」，不能塞進 AbstractHandler
}

func (oSelf *UserConsumer) AddUser(msg amqp.Delivery) {
	var payload createUserMessage
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		pkg.Logger(pkg.Consumer).Error("AddUser 訊息格式錯誤",
			zap.Error(err),
		)
		msg.Nack(false, false)
		return
	}

	if _, err := oSelf.userUsecase.AddUserByName(payload.Name); err != nil {
		pkg.Logger(pkg.Consumer).Error("AddUser 失敗",
			zap.String("name", payload.Name),
			zap.Error(err),
		)
		msg.Nack(false, true)
		return
	}

	pkg.Logger(pkg.Consumer).Info("AddUser 成功",
		zap.String("name", payload.Name),
	)
	msg.Ack(false)
}

func (oSelf *UserConsumer) Close() error {
	return oSelf.Conn.Close()
}
