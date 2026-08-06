//go:build wireinject
// +build wireinject

package container

import (
	"database/sql"

	"github.com/google/wire"
	"google.golang.org/grpc"

	"example/bootstrap"
	"example/helper"
	"example/input/abstract"
	inputClient "example/input/client"
	inputConsumer "example/input/consumer"
	inputGrpc "example/input/grpc"
	inputHttp "example/input/http"
	inputPort "example/input/port"
	outputMysql "example/output/mysql"
	"example/usecase"
)

/*
main.go 原本手動 new 出每一個物件，這裡改用 Wire 組裝。
main.go 裡的三個連線字串（mysql dsn / amqp dsn / grpc addr）沒有配置檔可讀，
所以做成三個獨立型別（MysqlDSN / AmqpDSN / GrpcAddr）當作 InitContainer 的參數，
讓 Wire 能依型別分辨要注入到哪個 provider，呼叫端則跟 main.go 原本的字串一致。
*/

type MysqlDSN string
type AmqpDSN string
type GrpcAddr string

// provideMysqlDB / provideGrpcConn / provideUserConsumer 只是把型別化的 DSN
// 轉成 bootstrap 原本要的 string，本身邏輯完全沒動。
func provideMysqlDB(sDSN MysqlDSN) (*sql.DB, error) {
	return bootstrap.NewMysql(string(sDSN))
}

func provideGrpcConn(sAddr GrpcAddr) (*grpc.ClientConn, error) {
	return bootstrap.NewClient(string(sAddr))
}

func provideUserConsumer(sDSN AmqpDSN, oUsecase inputPort.UserUsecase, oAbstractHandler *abstract.AbstractHandler) (*inputConsumer.UserConsumer, error) {
	return inputConsumer.NewUserConsumer(string(sDSN), oUsecase, oAbstractHandler)
}

type Container struct {
	// Helper
	*helper.AbstractHelper
	*helper.AesHelper
	*abstract.AbstractHandler

	// Input Adapter：四種輸入來源共用同一個 UserUsecase
	ConsumerUserHandler *inputConsumer.UserConsumer // MQ 消費者
	ClientUserHandler   *inputClient.UserHandler    // gRPC client stream 訂閱
	GrpcUserHandler     *inputGrpc.UserHandler      // gRPC server
	HttpUserHandler     *inputHttp.UserHandler      // HTTP handler
}

func InitContainer(sMysqlDSN MysqlDSN, sAmqpDSN AmqpDSN, sGrpcAddr GrpcAddr) (*Container, error) {
	wire.Build(
		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		abstract.NewAbstractHandler,

		// bootstrap
		provideMysqlDB,
		provideGrpcConn,

		// output
		outputMysql.NewUserRepository,

		// usecase
		usecase.NewUserUsecase,

		// input
		provideUserConsumer,
		inputClient.NewUserHandler,
		inputGrpc.NewUserHandler,
		inputHttp.NewUserHandler,

		wire.Struct(new(Container), "*"),
	)
	return nil, nil
}
