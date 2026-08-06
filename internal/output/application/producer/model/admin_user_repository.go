package producer

import (
	"encoding/json"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"

	domain "example/internal/domain"

	outputApplicationProducer "example/internal/output/application/producer"
)

type AdminUserRepository struct {
	*outputApplicationProducer.AbstractRepository
}

func NewAdminUserRepository(oAbstractRepository *outputApplicationProducer.AbstractRepository) (*AdminUserRepository, error) {
	_, err := oAbstractRepository.Channel.QueueDeclare("AdminUser.AddOne", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &AdminUserRepository{
		AbstractRepository: oAbstractRepository,
	}, nil
}

func (oSelf *AdminUserRepository) AddOne(oAdminUser *domain.AdminUser) error {
	body, err := json.Marshal(oAdminUser)
	if err != nil {
		return err
	}

	return oSelf.Channel.Publish(
		"",
		"AdminUser.AddOne",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (oSelf *AdminUserRepository) ShowOneById(id uint) (*domain.AdminUser, error) {
	return nil, errors.New("not supported by producer")
}

// Close 是空實作：Conn／Channel 現在都是從 AbstractRepository 注入的共用資源，
// 生命週期不屬於這個 repository，不該由這裡關閉。
func (oSelf *AdminUserRepository) Close() error {
	return nil
}
