//go:build wireinject
// +build wireinject

package container

/*
	usecase.NewUserUsecase的簽名是： func NewUserUsecase(oAbstractUsecase *AbstractUsecase) inputPort.UserUsecase
	回傳型別宣告的是介面 inputPort.UserUsecase，不是具體型別 *usecase.UserUsecase。

	wire 是純靜態分析工具，它只看 provider 函式簽名上寫的型別去做「型別對型別」的精確匹配，不會去看函式內部實際 return &UserUsecase{...} 塞的是什麼具體型別。所以 wire 註冊到的 provider 是「能生產 inputPort.UserUsecase」，而不是「能生產 *usecase.UserUsecase」——即使後者在執行期其實是同一個值。
*/

import (
	"github.com/google/wire"

	pkg "example/pkg"

	bootstrap "example/internal/bootstrap"
	internalClient "example/internal/client"

	helper "example/internal/helper"

	inputClient "example/internal/input/client"
	inputCommand "example/internal/input/command"
	inputConsumer "example/internal/input/consumer"
	inputCron "example/internal/input/cron"
	inputFacade "example/internal/input/facade"
	inputFacadeGame "example/internal/input/facade/game"
	inputFacadeRegister "example/internal/input/facade/register"
	inputFacadeTable "example/internal/input/facade/table"
	inputHttpAdmin "example/internal/input/http/admin"
	inputHttpAdminResource "example/internal/input/http/admin/resource"
	inputWebsocket "example/internal/input/websocket"

	Resource "example/internal/input/resource"
	ResourceModel "example/internal/input/resource/model"

	MiddlewareAdmin "example/internal/middleware/admin"

	inputPort "example/internal/usecase/port"

	usecase "example/internal/usecase"
	usecaseResource "example/internal/usecase/resource"

	cache "example/internal/output/cache"
	memory "example/internal/output/memory"
	mysql "example/internal/output/mysql"
)

/*

 */

type FacadeContainer struct {

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper
	*helper.RsaHelper
	*helper.LoggerHelper

	*usecase.AbstractUsecase
	inputPort.UserUsecase

	// gRPC Facade server
	FacadeAbstract           *inputFacade.AbstractHandler
	FacadeGameUser           *inputFacadeGame.UserHandler
	inputFacadeTableScanner  *inputFacadeTable.ScannerHandler
	FacadeTableAuthenticator *inputFacadeRegister.AuthenticatorHandler
}

func InitFacadeContainer() (*FacadeContainer, error) {
	wire.Build(

		// bootstrap
		bootstrap.NewMysql,
		bootstrap.NewRedis,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewRsaHelper,
		helper.NewCacheHelper,
		helper.NewLoggerHelper,

		// input-facade
		inputFacade.NewAbstractHandler,
		inputFacadeGame.NewUserHandler,
		inputFacadeTable.NewScannerHandler,
		inputFacadeRegister.NewAuthenticatorHandler,

		// usecase
		usecase.NewAbstractUsecase,
		usecase.NewUserUsecase,

		// output
		cache.NewUserRepository,

		wire.Struct(new(FacadeContainer), "*"),
	)
	return nil, nil
}

//////////////////////////////////////////////////////////////////////////////

type ResourceContainer struct {

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper
	*helper.RsaHelper
	*helper.LoggerHelper

	*usecase.AbstractUsecase
	inputPort.AdminUserUsecase

	// gRPC Resource server
	ResourceAbstract       *Resource.AbstractHandler
	ResourceModelAdminUser *ResourceModel.AdminUserHandler
}

func InitResourceContainer() (*ResourceContainer, error) {
	wire.Build(

		// bootstrap
		bootstrap.NewMysql,
		bootstrap.NewRedis,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewRsaHelper,
		helper.NewCacheHelper,
		helper.NewLoggerHelper,

		// input-resource
		Resource.NewAbstractHandler,
		ResourceModel.NewAdminUserHandler,

		// usecase
		usecase.NewAbstractUsecase,
		usecaseResource.NewAbstractUsecase,
		usecaseResource.NewAdminUserUsecase,

		// output
		cache.NewUserRepository,
		mysql.NewAdminUserRepository,

		wire.Struct(new(ResourceContainer), "*"),
	)
	return nil, nil
}

// HttpContainer 只給 `http` Gin 服務使用。
type HttpContainer struct {

	// pkg
	*pkg.Response

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper
	*helper.RsaHelper
	*helper.LoggerHelper

	*usecase.AbstractUsecase
	inputPort.UserUsecase

	// HTTP server -Controller
	HttpAdminResourceUser *inputHttpAdminResource.UserHandler

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
}

func InitHttpContainer() (*HttpContainer, error) {
	wire.Build(

		// pkg
		pkg.NewResponse,

		// bootstrap
		bootstrap.NewMysql,
		bootstrap.NewRedis,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewRsaHelper,
		helper.NewCacheHelper,
		helper.NewLoggerHelper,

		// input-http
		inputHttpAdmin.NewAbstractHandler,
		inputHttpAdminResource.NewUserHandler,

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
		cache.NewUserRepository,

		wire.Struct(new(HttpContainer), "*"),
	)
	return nil, nil
}

//////////////////////////////////////////////////////////////////////////////

// ConsumerContainer 只給 `consumer` MQ 消費者服務使用。
type ConsumerContainer struct {

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper

	*usecase.AbstractUsecase
	inputPort.UserUsecase

	// MQ 消費者
	ConsumerUser *inputConsumer.UserConsumer
}

func InitConsumerContainer() (*ConsumerContainer, error) {
	wire.Build(

		// bootstrap
		bootstrap.NewMysql,
		bootstrap.NewAmqp,
		bootstrap.NewRedis,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewCacheHelper,

		// input-consumer
		inputConsumer.NewAbstractHandler,
		inputConsumer.NewUserConsumer,

		// usecase
		usecase.NewAbstractUsecase,
		usecase.NewUserUsecase,

		// output
		cache.NewUserRepository,

		wire.Struct(new(ConsumerContainer), "*"),
	)
	return nil, nil
}

//////////////////////////////////////////////////////////////////////////////

// CronContainer 只給 `cron` 排程服務使用。
type CronContainer struct {

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper

	*usecase.AbstractUsecase
	inputPort.UserUsecase

	// 排程 server
	CronUser *inputCron.UserCron
}

func InitCronContainer() (*CronContainer, error) {
	wire.Build(

		// bootstrap
		bootstrap.NewMysql,
		bootstrap.NewRedis,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewCacheHelper,

		// input-cron
		inputCron.NewAbstractHandler,
		inputCron.NewUserCron,

		// usecase
		usecase.NewAbstractUsecase,
		usecase.NewUserUsecase,

		// output
		cache.NewUserRepository,

		wire.Struct(new(CronContainer), "*"),
	)
	return nil, nil
}

//////////////////////////////////////////////////////////////////////////////

// WebsocketContainer 只給 `websocket` 服務使用。
type WebsocketContainer struct {

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper

	*usecase.AbstractUsecase
	inputPort.UserUsecase

	// websocket server
	WebsocketUser *inputWebsocket.UserHandler
}

func InitWebsocketContainer() (*WebsocketContainer, error) {
	wire.Build(

		// bootstrap
		bootstrap.NewMysql,
		bootstrap.NewRedis,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewCacheHelper,

		// input-websocket
		inputWebsocket.NewAbstractHandler,
		inputWebsocket.NewUserHandler,

		// usecase
		usecase.NewAbstractUsecase,
		usecase.NewUserUsecase,

		// output
		cache.NewUserRepository,

		wire.Struct(new(WebsocketContainer), "*"),
	)
	return nil, nil
}

//////////////////////////////////////////////////////////////////////////////

// ClientContainer 只給 `client` （訂閱外部 gRPC stream）服務使用。
type ClientContainer struct {

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper

	*usecase.AbstractUsecase
	inputPort.UserUsecase

	// gRPC client stream 訂閱
	ClientUser *inputClient.UserHandler
}

func InitClientContainer() (*ClientContainer, error) {
	wire.Build(

		// bootstrap
		bootstrap.NewMysql,
		bootstrap.NewClient,
		bootstrap.NewRedis,

		//
		internalClient.NewClient,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewCacheHelper,

		// input-client
		inputClient.NewAbstractHandler,
		inputClient.NewUserHandler,

		// usecase
		usecase.NewAbstractUsecase,
		usecase.NewUserUsecase,

		// output
		cache.NewUserRepository,

		wire.Struct(new(ClientContainer), "*"),
	)
	return nil, nil
}

// 、
type CommandContainer struct {

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper

	// command
	*inputCommand.AbstractHandler
	*inputCommand.UserHandler

	*usecase.AbstractUsecase
	inputPort.UserUsecase
}

func InitCommandContainer() (*CommandContainer, error) {
	wire.Build(

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,

		inputCommand.NewAbstractHandler,
		inputCommand.NewUserHandler,

		// usecase
		usecase.NewAbstractUsecase,
		usecase.NewUserUsecase,

		// output
		memory.NewUserRepository,

		wire.Struct(new(CommandContainer), "*"),
	)
	return nil, nil
}
