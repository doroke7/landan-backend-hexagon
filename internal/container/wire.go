//go:build wireinject
// +build wireinject

package container

import (
	"github.com/google/wire"

	pkg "example/pkg"

	bootstrap "example/internal/bootstrap"
	internalClient "example/internal/client"
	helper "example/internal/helper"
	client "example/internal/input/client"
	consumer "example/internal/input/consumer"
	cron "example/internal/input/cron"
	Facade "example/internal/input/facade"
	FacadeGame "example/internal/input/facade/game"
	FacadeTable "example/internal/input/facade/table"

	MiddlewareAdmin "example/internal/middleware/admin"

	HttpAdmin "example/internal/input/http/admin"
	HttpAdminResource "example/internal/input/http/admin/resource"

	inputPort "example/internal/input/port"
	"example/internal/input/websocket"
	"example/internal/output/cache"
	"example/internal/output/mysql"
	"example/internal/usecase"
)

/*

 */

type Container struct {

	// pkg
	*pkg.Response

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper
	*helper.RsaHelper
	*helper.LoggerHelper

	*usecase.AbstractUsecase
	inputPort.UserUsecase
	/*
		usecase.NewUserUsecase的簽名是： func NewUserUsecase(oAbstractUsecase *AbstractUsecase) inputPort.UserUsecase
		回傳型別宣告的是介面 inputPort.UserUsecase，不是具體型別 *usecase.UserUsecase。

		wire 是純靜態分析工具，它只看 provider 函式簽名上寫的型別去做「型別對型別」的精確匹配，不會去看函式內部實際 return &UserUsecase{...} 塞的是什麼具體型別。所以 wire 註冊到的 provider 是「能生產 inputPort.UserUsecase」，而不是「能生產 *usecase.UserUsecase」——即使後者在執行期其實是同一個值。
	*/

	// Input Adapter：四種輸入來源共用同一個 UserUsecase

	// MQ 消費者
	ConsumerUser *consumer.UserConsumer

	// gRPC client stream 訂閱
	ClientUser *client.UserHandler

	// gRPC server
	FacadeGameUser         *FacadeGame.UserHandler
	FacadeTableScannerUser *FacadeTable.ScannerHandler

	// HTTP server -Controller
	HttpAdminResourceUser *HttpAdminResource.UserHandler

	// HTTP server -Middleware
	// Middleware 部分
	AdminAbstractMiddleware       *MiddlewareAdmin.AbstractMiddleware
	AdminAdminMiddleware          *MiddlewareAdmin.AdminMiddleware
	AdminAuthenticationMiddleware *MiddlewareAdmin.AuthenticationMiddleware
	AdminDecryptionMiddleware     *MiddlewareAdmin.DecryptionMiddleware
	AdminEncryptionMiddleware     *MiddlewareAdmin.EncryptionMiddleware
	AdminErrorMiddleware          *MiddlewareAdmin.ErrorMiddleware
	AdminLoggerMiddleware         *MiddlewareAdmin.LoggerMiddleware
	AdminNonexistentMiddleware    *MiddlewareAdmin.NonexistentMiddleware
	AdminRequestMiddleware        *MiddlewareAdmin.RequestMiddleware
	AdminResponseMiddleware       *MiddlewareAdmin.ResponseMiddleware
	AdminSignatureMiddleware      *MiddlewareAdmin.SignatureMiddleware
	// 排程 server
	CronUser *cron.UserCron

	// websocket server
	WebsocketUser *websocket.UserHandler
}

func InitContainer() (*Container, error) {
	wire.Build(

		// pkg
		pkg.NewResponse,

		// bootstrap
		bootstrap.NewMysql,
		bootstrap.NewClient,
		bootstrap.NewAmqp,
		bootstrap.NewRedis,

		//
		internalClient.NewClient,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewRsaHelper,
		helper.NewCacheHelper,
		helper.NewLoggerHelper,

		// input-consumer
		consumer.NewAbstractHandler,
		consumer.NewUserConsumer,

		// input-http
		HttpAdmin.NewAbstractHandler,
		HttpAdminResource.NewUserHandler,

		// input-cron
		cron.NewAbstractHandler,
		cron.NewUserCron,

		// input-websocket
		websocket.NewAbstractHandler,
		websocket.NewUserHandler,

		// input-facade
		Facade.NewAbstractHandler,
		FacadeGame.NewUserHandler,
		FacadeTable.NewScannerHandler,

		// input-client
		client.NewAbstractHandler,
		client.NewUserHandler,

		// Middleware 部分
		MiddlewareAdmin.NewAbstractMiddleware,
		MiddlewareAdmin.NewAdminMiddleware,
		MiddlewareAdmin.NewAuthenticationMiddleware,
		MiddlewareAdmin.NewDecryptionMiddleware,
		MiddlewareAdmin.NewEncryptionMiddleware,
		MiddlewareAdmin.NewErrorMiddleware,
		MiddlewareAdmin.NewLoggerMiddleware,
		MiddlewareAdmin.NewNonexistentMiddleware,
		MiddlewareAdmin.NewRequestMiddleware,
		MiddlewareAdmin.NewResponseMiddleware,
		MiddlewareAdmin.NewSignatureMiddleware,

		// usecase
		usecase.NewAbstractUsecase,
		usecase.NewUserUsecase,

		// output
		mysql.NewUserRepository,
		cache.NewUserRepository,

		wire.Struct(new(Container), "*"),
	)
	return nil, nil
}
