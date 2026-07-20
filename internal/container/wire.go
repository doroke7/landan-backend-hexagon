//go:build wireinject
// +build wireinject

package container

/*

 */

import (
	"github.com/google/wire"

	pkg "example/pkg"

	bootstrap "example/bootstrap"
	Client "example/internal/client"

	helper "example/internal/helper"

	MiddlewareAdmin "example/internal/middleware/admin"

	InterceptorFacadeAdmin "example/internal/interceptor/facade/game"
	InterceptorResource "example/internal/interceptor/resource"

	inputCommand "example/internal/input/application/command/admin/resource"
	inputConsumer "example/internal/input/application/consumer/admin/resource"
	inputCron "example/internal/input/application/cron/admin/resource"
	inputFacade "example/internal/input/application/facade"
	inputFacadeRegister "example/internal/input/application/facade/register"
	inputFacadeTable "example/internal/input/application/facade/table"
	inputHttpAdmin "example/internal/input/application/http/admin"
	inputHttpAdminAuthentication "example/internal/input/application/http/admin/authentication"

	inputResource "example/internal/input/application/resource"
	inputResourceModel "example/internal/input/application/resource/model"

	usecasePortResourceModel "example/internal/usecase/port/resource/model"

	usecaseApplicationCommand "example/internal/usecase/application/command/admin/resource"
	usecaseApplicationConsumer "example/internal/usecase/application/consumer/admin/resource"
	usecaseApplicationCron "example/internal/usecase/application/cron/admin/resource"
	usecaseApplicationHttpAdminAuthentication "example/internal/usecase/application/http/admin/authentication"
	usecaseApplicationResourceModel "example/internal/usecase/application/resource/model"

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
	*helper.JwtHelper

	// Clients
	ResourceClient *Client.ResourceClient

	// HTTP server -Controller
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
		bootstrap.NewResource,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewRsaHelper,
		helper.NewJwtHelper,

		// output
		outputApplicationResourceModel.NewAdminUserRepository,

		// usecase
		usecaseApplicationHttpAdminAuthentication.NewAbstractUsecase,
		usecaseApplicationHttpAdminAuthentication.NewAuthenticatorUsecase,

		// client
		Client.NewModel,
		Client.NewResourceClient,

		// input-http
		inputHttpAdmin.NewAbstractHandler,

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

	// gRPC Facade server
	FacadeAbstract           *inputFacade.AbstractHandler
	FacadeTableScanner       *inputFacadeTable.ScannerHandler
	FacadeTableAuthenticator *inputFacadeRegister.AuthenticatorHandler

	// gRPC Facade Interceptor
	FacadeAdminErrorInterceptor          *InterceptorFacadeAdmin.ErrorInterceptor
	FacadeAdminStatusInterceptor         *InterceptorFacadeAdmin.StatusInterceptor
	FacadeAdminLoggerInterceptor         *InterceptorFacadeAdmin.LoggerInterceptor
	FacadeAdminAuthenticationInterceptor *InterceptorFacadeAdmin.AuthenticationInterceptor
}

func InitFacadeContainer() (*FacadeContainer, error) {
	wire.Build(

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewRsaHelper,

		// input-facade
		inputFacade.NewAbstractHandler,
		inputFacadeTable.NewScannerHandler,
		inputFacadeRegister.NewAuthenticatorHandler,

		// interceptor-facade
		InterceptorFacadeAdmin.NewAbstractInterceptor,
		InterceptorFacadeAdmin.NewErrorInterceptor,
		InterceptorFacadeAdmin.NewStatusInterceptor,
		InterceptorFacadeAdmin.NewLoggerInterceptor,
		InterceptorFacadeAdmin.NewAuthenticationInterceptor,

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

	// gRPC Resource Interceptor
	ResourceAllInterceptor *InterceptorResource.AllInterceptor
}

func InitResourceContainer() (*ResourceContainer, error) {
	wire.Build(

		// bootstrap
		bootstrap.NewMysql,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewRsaHelper,

		// output
		outputApplicationMysqlModel.NewAbstractRepository,
		outputApplicationMysqlModel.NewAdminUserRepository,

		// usecase
		usecaseApplicationResourceModel.NewAbstractUsecase,
		usecaseApplicationResourceModel.NewAdminUserUsecase,

		// input-resource
		inputResource.NewAbstractHandler,
		inputResourceModel.NewAdminUserHandler,

		// interceptor-resource
		InterceptorResource.NewAbstractInterceptor,
		InterceptorResource.NewAllInterceptor,

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

	// MQ 消費者
	*inputConsumer.AbstractHandler
	ConsumerAppUser *inputConsumer.AppUserHandler
}

func InitConsumerContainer() (*ConsumerContainer, error) {
	wire.Build(

		// bootstrap
		bootstrap.NewMysql,
		bootstrap.NewAmqp,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,

		// output
		outputApplicationMysqlModel.NewAbstractRepository,
		outputApplicationMysqlModel.NewAppUserRepository,

		// usecase
		usecaseApplicationConsumer.NewAppUserUsecase,

		// input-consumer
		inputConsumer.NewAbstractHandler,
		inputConsumer.NewAppUserHandler,

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

	// 排程 server
	CronAppUser *inputCron.AppUserHandler
}

func InitCronContainer() (*CronContainer, error) {
	wire.Build(

		// bootstrap
		bootstrap.NewMysql,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,

		// output
		outputApplicationMysqlModel.NewAbstractRepository,
		outputApplicationMysqlModel.NewAppUserRepository,

		// usecase
		usecaseApplicationCron.NewAppUserUsecase,

		// input-cron
		inputCron.NewAbstractHandler,
		inputCron.NewAppUserHandler,

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
}

func InitWebsocketContainer() (*WebsocketContainer, error) {
	wire.Build(

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,

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
}

func InitClientContainer() (*ClientContainer, error) {
	wire.Build(

		// helper 部份
		helper.NewAbstractHelper,
		helper.NewAesHelper,

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
	CommandAppUser *inputCommand.AppUserHandler
}

func InitCommandContainer() (*CommandContainer, error) {
	wire.Build(

		// bootstrap
		bootstrap.NewMysql,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,

		// output
		outputApplicationMysqlModel.NewAbstractRepository,
		outputApplicationMysqlModel.NewAppUserRepository,

		// usecase
		usecaseApplicationCommand.NewAppUserUsecase,

		// command
		inputCommand.NewAbstractHandler,
		inputCommand.NewAppUserHandler,

		wire.Struct(new(CommandContainer), "*"),
	)
	return nil, nil
}
