//go:build wireinject
// +build wireinject

package container

/*

 */

import (
	"github.com/google/wire"

	pkg "example/pkg"

	bootstrap "example/internal/bootstrap"
	Client "example/internal/client"

	helper "example/internal/helper"

	MiddlewareAdmin "example/internal/middleware/admin"

	inputClient "example/internal/input/application/client"
	inputCommand "example/internal/input/application/command"
	inputConsumer "example/internal/input/application/consumer"
	inputCron "example/internal/input/application/cron"
	inputFacade "example/internal/input/application/facade"
	inputFacadeGame "example/internal/input/application/facade/game"
	inputFacadeRegister "example/internal/input/application/facade/register"
	inputFacadeTable "example/internal/input/application/facade/table"
	inputHttpAdmin "example/internal/input/application/http/admin"
	inputHttpAdminAuthentication "example/internal/input/application/http/admin/authentication"
	inputHttpAdminResource "example/internal/input/application/http/admin/resource"
	inputWebsocket "example/internal/input/application/websocket"

	inputResource "example/internal/input/application/resource"
	inputResourceModel "example/internal/input/application/resource/model"

	usecasePortFacadeModel "example/internal/usecase/port/facade/model"
	usecasePortResourceModel "example/internal/usecase/port/resource/model"

	usecaseApplicationFacadeModel "example/internal/usecase/application/facade/model"
	usecaseApplicationResourceModel "example/internal/usecase/application/resource/model"

	outputApplicationCacheModel "example/internal/output/application/cache/model"
	outputApplicationMemoryModel "example/internal/output/application/memory/model"
	outputApplicationMysqlModel "example/internal/output/application/mysql/model"
	outputApplicationResourceModel "example/internal/output/application/resource/model"
)

// HttpContainer 只給 `http` Gin 服務使用。
type HttpContainer struct {

	// pkg
	*pkg.Response

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper
	*helper.RsaHelper

	*usecaseApplicationFacadeModel.AbstractUsecase
	usecasePortFacadeModel.UserUsecase

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
		usecaseApplicationFacadeModel.NewAbstractUsecase,
		usecaseApplicationFacadeModel.NewUserUsecase,

		// output
		outputApplicationCacheModel.NewUserRepository,
		outputApplicationResourceModel.NewAdminUserRepository,

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

	*usecaseApplicationFacadeModel.AbstractUsecase
	usecasePortFacadeModel.UserUsecase

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
		usecaseApplicationFacadeModel.NewAbstractUsecase,
		usecaseApplicationFacadeModel.NewUserUsecase,

		// output
		outputApplicationCacheModel.NewUserRepository,

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

	*usecaseApplicationResourceModel.AbstractUsecase
	usecasePortResourceModel.AdminUserUsecase

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
		usecaseApplicationResourceModel.NewAbstractUsecase,
		usecaseApplicationResourceModel.NewAdminUserUsecase,

		// output
		outputApplicationMysqlModel.NewAdminUserRepository,

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

	*usecaseApplicationFacadeModel.AbstractUsecase
	usecasePortFacadeModel.UserUsecase

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
		usecaseApplicationFacadeModel.NewAbstractUsecase,
		usecaseApplicationFacadeModel.NewUserUsecase,

		// output
		outputApplicationCacheModel.NewUserRepository,

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

	*usecaseApplicationFacadeModel.AbstractUsecase
	usecasePortFacadeModel.UserUsecase

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
		usecaseApplicationFacadeModel.NewAbstractUsecase,
		usecaseApplicationFacadeModel.NewUserUsecase,

		// output
		outputApplicationCacheModel.NewUserRepository,

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

	*usecaseApplicationFacadeModel.AbstractUsecase
	usecasePortFacadeModel.UserUsecase

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
		usecaseApplicationFacadeModel.NewAbstractUsecase,
		usecaseApplicationFacadeModel.NewUserUsecase,

		// output
		outputApplicationCacheModel.NewUserRepository,

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

	*usecaseApplicationFacadeModel.AbstractUsecase
	usecasePortFacadeModel.UserUsecase

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
		usecaseApplicationFacadeModel.NewAbstractUsecase,
		usecaseApplicationFacadeModel.NewUserUsecase,

		// output
		outputApplicationCacheModel.NewUserRepository,

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

	*usecaseApplicationFacadeModel.AbstractUsecase
	// usecasePortFacadeModel.UserUsecase
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
		usecaseApplicationFacadeModel.NewAbstractUsecase,
		usecaseApplicationFacadeModel.NewUserUsecase,

		// output
		outputApplicationMemoryModel.NewUserRepository,

		wire.Struct(new(CommandContainer), "*"),
	)
	return nil, nil
}
