package producer

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractRepository 做法一致。
// Conn／Channel 是從 bootstrap 注入的共用 amqp 連線與 channel，
// 跟 mysql.AbstractRepository 持有 *gorm.DB 是同一種角色。
type AbstractRepository struct {
	Context context.Context
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewAbstractRepository(oContext context.Context, oConn *amqp.Connection) (*AbstractRepository, error) {
	oChannel, err := oConn.Channel()
	if err != nil {
		return nil, err
	}

	return &AbstractRepository{
		Context: oContext,
		Conn:    oConn,
		Channel: oChannel,
	}, nil
}
