package producer

import (
	"encoding/json"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"

	"example/internal/domain"
	"example/internal/output/port/any/model"
)

type UserProducer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewUserProducer(dsn string) (port.UserRepository, error) {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare("user.created", true, false, false, false, nil)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &UserProducer{conn: conn, channel: ch}, nil
}

func (oSelf *UserProducer) AddOne(user *domain.User) error {
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return oSelf.channel.Publish(
		"",
		"user.created",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (oSelf *UserProducer) ShowOneById(id int) (*domain.User, error) {
	return nil, errors.New("not supported by producer")
}

func (oSelf *UserProducer) Close() error {
	oSelf.channel.Close()
	return oSelf.conn.Close()
}
