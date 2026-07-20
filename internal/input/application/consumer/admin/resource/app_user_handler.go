package consumer

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"

	inputConsumer "example/internal/input/application/consumer"
	port "example/internal/usecase/port/consumer/admin/resource"
	pkg "example/pkg"

	"go.uber.org/zap"
)

type increaseBalanceMessage struct {
	Id     uint `json:"id"`
	Amount uint `json:"amount"`
}

type AppUserHandler struct {
	*inputConsumer.AbstractHandler
	appUserUsecase port.AppUserUsecase
}

func NewAppUserHandler(oAppUserUsecase port.AppUserUsecase, oAbstractHandler *inputConsumer.AbstractHandler) *AppUserHandler {
	return &AppUserHandler{
		AbstractHandler: oAbstractHandler,
		appUserUsecase:  oAppUserUsecase,
	}
}

func (oSelf *AppUserHandler) IncreaseBalance(msg amqp.Delivery) {
	var payload increaseBalanceMessage
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		pkg.Logger(pkg.Consumer).Error("IncreaseBalance 訊息格式錯誤",
			zap.Error(err),
		)
		msg.Nack(false, false)
		return
	}

	oAppUser, err := oSelf.appUserUsecase.IncreaseBalance(payload.Id, payload.Amount)
	if err != nil {
		pkg.Logger(pkg.Consumer).Error("IncreaseBalance 失敗",
			zap.Uint("id", payload.Id),
			zap.Error(err),
		)
		msg.Nack(false, true)
		return
	}

	pkg.Logger(pkg.Consumer).Info("IncreaseBalance 成功",
		zap.Uint("id", oAppUser.Id),
		zap.Uint("balance", oAppUser.Balance),
	)
	msg.Ack(false)
}
