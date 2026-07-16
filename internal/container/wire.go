//go:build wireinject
// +build wireinject

package container

/*
	usecaseFacadeModelApplication.NewUserUsecase的簽名是： func NewUserUsecase(oAbstractUsecase *AbstractUsecase) usecaseFacadeModelPort.UserUsecase
	回傳型別宣告的是介面 usecaseFacadeModelPort.UserUsecase，不是具體型別 *usecaseFacadeModelApplication.UserUsecase。

	wire 是純靜態分析工具，它只看 provider 函式簽名上寫的型別去做「型別對型別」的精確匹配，不會去看函式內部實際 return &UserUsecase{...} 塞的是什麼具體型別。所以 wire 註冊到的 provider 是「能生產 usecaseFacadeModelPort.UserUsecase」，而不是「能生產 *usecaseFacadeModelApplication.UserUsecase」——即使後者在執行期其實是同一個值。

	AdminUserUsecase（resource 服務專屬）走的是 usecaseResourceModelApplication / usecaseResourceModelPort，
	跟 facade 以及其他週邊 adapter 共用的 UserUsecase 是完全獨立的兩個 package，
	因為兩邊的商業邏輯差異太大，故意不共用同一個 AbstractUsecase。
*/

import (
	"github.com/google/wire"

	pkg "example/pkg"

	bootstrap "example/internal/bootstrap"
	Client "example/internal/client"

	helper "example/internal/helper"

	MiddlewareAdmin "example/internal/middleware/admin"

	inputClient "example/internal/input/client"
	inputCommand "example/internal/input/command"
	inputConsumer "example/internal/input/consumer"
	inputCron "example/internal/input/cron"
	inputFacade "example/internal/input/facade"
	inputFacadeGame "example/internal/input/facade/game"
	inputFacadeRegister "example/internal/input/facade/register"
	inputFacadeTable "example/internal/input/facade/table"
	inputHttpAdmin "example/internal/input/http/admin"
	inputHttpAdminAuthentication "example/internal/input/http/admin/authentication"
	inputHttpAdminResource "example/internal/input/http/admin/resource"
	inputWebsocket "example/internal/input/websocket"

	inputResource "example/internal/input/resource"
	inputResourceModel "example/internal/input/resource/model"

	usecaseFacadeModelPort "example/internal/usecase/facade/model/port"
	usecaseResourceModelPort "example/internal/usecase/resource/model/port"

	usecaseFacadeModelApplication "example/internal/usecase/facade/model/application"
	usecaseResourceModelApplication "example/internal/usecase/resource/model/application"

	outputCache "example/internal/output/cache/model"
	outputMemory "example/internal/output/memory/model"
	outputMysql "example/internal/output/mysql/model"
	outputResourceModel "example/internal/output/resource/model"
)

// HttpContainer 只給 `http` Gin 服務使用。
type HttpContainer struct {

	// pkg
	*pkg.Response

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper
	*helper.RsaHelper

	*usecaseFacadeModelApplication.AbstractUsecase
	usecaseFacadeModelPort.UserUsecase

	// Clients
	ResourceClient *Client.ResourceClient

	// HTTP server -Controller
	HttpAdminResourceUser                *inputHttpAdminResource.UserHandler
	HttpAdminAuthenticationAuthenticator *inputHttpAdminAuthentication.AuthenticatorHandler

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
		bootstrap.NewResource,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewRsaHelper,
		helper.NewCacheHelper,

		// usecase
		usecaseFacadeModelApplication.NewAbstractUsecase,
		usecaseFacadeModelApplication.NewUserUsecase,

		// output
		outputCache.NewUserRepository,
		outputResourceModel.NewAdminUserRepository,

		// client
		Client.NewModel,
		Client.NewResourceClient,

		// input-http
		inputHttpAdmin.NewAbstractHandler,
		inputHttpAdminResource.NewUserHandler,

		inputHttpAdminAuthentication.NewAuthenticatorHandler,

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

		wire.Struct(new(HttpContainer), "*"),
	)
	return nil, nil
}

type FacadeContainer struct {

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper
	*helper.RsaHelper

	*usecaseFacadeModelApplication.AbstractUsecase
	usecaseFacadeModelPort.UserUsecase

	// gRPC Facade server
	FacadeAbstract           *inputFacade.AbstractHandler
	FacadeGameUser           *inputFacadeGame.UserHandler
	FacadeTableScanner       *inputFacadeTable.ScannerHandler
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

		// input-facade
		inputFacade.NewAbstractHandler,
		inputFacadeGame.NewUserHandler,
		inputFacadeTable.NewScannerHandler,
		inputFacadeRegister.NewAuthenticatorHandler,

		// usecase
		usecaseFacadeModelApplication.NewAbstractUsecase,
		usecaseFacadeModelApplication.NewUserUsecase,

		// output
		outputCache.NewUserRepository,

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

	*usecaseResourceModelApplication.AbstractUsecase
	usecaseResourceModelPort.AdminUserUsecase

	// gRPC Resource server
	ResourceAbstract       *inputResource.AbstractHandler
	ResourceModelAdminUser *inputResourceModel.AdminUserHandler
}

func InitResourceContainer() (*ResourceContainer, error) {
	wire.Build(

		// bootstrap
		bootstrap.NewMysql,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewRsaHelper,

		// input-resource
		inputResource.NewAbstractHandler,
		inputResourceModel.NewAdminUserHandler,

		// usecase
		usecaseResourceModelApplication.NewAbstractUsecase,
		usecaseResourceModelApplication.NewAdminUserUsecase,

		// output
		outputMysql.NewAdminUserRepository,

		wire.Struct(new(ResourceContainer), "*"),
	)
	return nil, nil
}

//////////////////////////////////////////////////////////////////////////////

// ConsumerContainer 只給 `consumer` MQ 消費者服務使用。
type ConsumerContainer struct {

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper

	*usecaseFacadeModelApplication.AbstractUsecase
	usecaseFacadeModelPort.UserUsecase

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
		usecaseFacadeModelApplication.NewAbstractUsecase,
		usecaseFacadeModelApplication.NewUserUsecase,

		// output
		outputCache.NewUserRepository,

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

	*usecaseFacadeModelApplication.AbstractUsecase
	usecaseFacadeModelPort.UserUsecase

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
		usecaseFacadeModelApplication.NewAbstractUsecase,
		usecaseFacadeModelApplication.NewUserUsecase,

		// output
		outputCache.NewUserRepository,

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

	*usecaseFacadeModelApplication.AbstractUsecase
	usecaseFacadeModelPort.UserUsecase

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
		usecaseFacadeModelApplication.NewAbstractUsecase,
		usecaseFacadeModelApplication.NewUserUsecase,

		// output
		outputCache.NewUserRepository,

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

	*usecaseFacadeModelApplication.AbstractUsecase
	usecaseFacadeModelPort.UserUsecase

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
		Client.NewClient,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewCacheHelper,

		// input-client
		inputClient.NewAbstractHandler,
		inputClient.NewUserHandler,

		// usecase
		usecaseFacadeModelApplication.NewAbstractUsecase,
		usecaseFacadeModelApplication.NewUserUsecase,

		// output
		outputCache.NewUserRepository,

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

	*usecaseFacadeModelApplication.AbstractUsecase
	// usecaseFacadeModelPort.UserUsecase
}

func InitCommandContainer() (*CommandContainer, error) {
	wire.Build(

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,

		// command
		inputCommand.NewAbstractHandler,
		inputCommand.NewUserHandler,

		// usecase
		usecaseFacadeModelApplication.NewAbstractUsecase,
		usecaseFacadeModelApplication.NewUserUsecase,

		// output
		outputMemory.NewUserRepository,

		wire.Struct(new(CommandContainer), "*"),
	)
	return nil, nil
}
