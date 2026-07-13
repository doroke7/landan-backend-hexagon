package consumer

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"example/internal/usecase/port"
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
		msg.Nack(false, false)
		return
	}

	if _, err := oSelf.userUsecase.AddUserByName(payload.Name); err != nil {
		log.Printf("create user failed: %v", err)
		msg.Nack(false, true)
		return
	}

	msg.Ack(false)
}

func (oSelf *UserConsumer) Close() error {
	return oSelf.Conn.Close()
}
