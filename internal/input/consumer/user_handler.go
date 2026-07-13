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

type UserConsumer struct {
	*AbstractHandler
	userUsecase port.UserUsecase // UserUsecase 是 driving port，是每個 handler 各自要注入的業務依賴，不是「抽象共用的技術基礎設施」，不能塞進 AbstractHandler
}

func NewUserConsumer(usecase port.UserUsecase, oAbstractHandler *AbstractHandler) (*UserConsumer, error) {
	return &UserConsumer{AbstractHandler: oAbstractHandler, userUsecase: usecase}, nil
}

// AddUser 是 "user.create" queue 對應的處理方法，職責跟 http/grpc handler 的
// AddUser 一樣單純：只管收到的訊息怎麼轉呼叫 usecase，佇列宣告／消費迴圈都交給 register.ConsumerRouter。
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
