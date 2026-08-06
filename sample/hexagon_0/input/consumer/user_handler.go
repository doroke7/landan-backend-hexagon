package consumer

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"example/input/abstract"
	"example/input/port"
)

type createUserMessage struct {
	Name string `json:"name"`
}

type UserConsumer struct {
	*abstract.AbstractHandler
	conn        *amqp.Connection
	userUsecase port.UserUsecase
}

func NewUserConsumer(dsn string, usecase port.UserUsecase, oAbstractHandler *abstract.AbstractHandler) (*UserConsumer, error) {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}
	return &UserConsumer{AbstractHandler: oAbstractHandler, conn: conn, userUsecase: usecase}, nil
}

func (oSelf *UserConsumer) Start(ctx context.Context) error {
	oChannel, err := oSelf.conn.Channel()
	if err != nil {
		return err
	}
	defer oChannel.Close()

	q, err := oChannel.QueueDeclare("user.create", true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := oChannel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return nil
			}
			oSelf.handle(msg)
		}
	}
}

func (oSelf *UserConsumer) handle(msg amqp.Delivery) {
	var payload createUserMessage
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		msg.Nack(false, false)
		return
	}

	if _, err := oSelf.userUsecase.CreateUser(payload.Name); err != nil {
		log.Printf("create user failed: %v", err)
		msg.Nack(false, true)
		return
	}

	msg.Ack(false)
}

func (oSelf *UserConsumer) Close() error {
	return oSelf.conn.Close()
}
