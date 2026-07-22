//go:build wireinject
// +build wireinject

package container

/*

 */

import (
	"github.com/google/wire"

	bootstrap "example/bootstrap"
	pkg "example/pkg"

	client "example/internal/client"

	helper "example/internal/helper"

	outputApplicationCache "example/internal/output/application/cache"
	outputApplicationCacheModel "example/internal/output/application/cache/model"

	outputApplicationMemory "example/internal/output/application/memory"
	outputApplicationMemoryModel "example/internal/output/application/memory/model"

	outputApplicationMysql "example/internal/output/application/mysql"
	outputApplicationMysqlModel "example/internal/output/application/mysql/model"
	outputApplicationResourceModel "example/internal/output/application/resource/model"

	usecasePortAnyModel "example/internal/usecase/port/any/model"

	usecaseApplicationAnyAdminAuthentication "example/internal/usecase/application/any/admin/authentication"
	usecaseApplicationAnyAdminResource "example/internal/usecase/application/any/admin/resource"
	usecaseApplicationAnyAnnouncement "example/internal/usecase/application/any/annoucement"
	usecaseApplicationAnyModel "example/internal/usecase/application/any/model"
	usecaseApplicationAnyWatcherSource "example/internal/usecase/application/any/watcher/source"

	middlewareAdmin "example/internal/middleware/admin"

	interceptorFacadeGame "example/internal/interceptor/facade/game"
	interceptorResource "example/internal/interceptor/resource"

	inputApplicationCommand "example/internal/input/application/command"
	inputApplicationCommandAdminAuthentication "example/internal/input/application/command/admin/authentication"
	inputApplicationCommandAdminResource "example/internal/input/application/command/admin/resource"

	inputApplicationConsumer "example/internal/input/application/consumer"
	inputApplicationConsumerAdminResource "example/internal/input/application/consumer/admin/resource"

	inputApplicationCron "example/internal/input/application/cron"
	inputApplicationCronAdminAuthentication "example/internal/input/application/cron/admin/authentication"
	inputApplicationCronAdminResource "example/internal/input/application/cron/admin/resource"

	inputApplicationSource "example/internal/input/application/source"
	inputApplicationSourceAnnouncement "example/internal/input/application/source/announcement"

	inputApplicationDaemon "example/internal/input/application/daemon"
	inputApplicationDaemonWatcherSource "example/internal/input/application/daemon/watcher/source"

	inputApplicationResource "example/internal/input/application/resource"
	inputApplicationResourceModel "example/internal/input/application/resource/model"

	inputApplicationFacade "example/internal/input/application/facade"
	inputApplicationFacadeRegister "example/internal/input/application/facade/register"
	inputApplicationFacadeTable "example/internal/input/application/facade/table"

	inputApplicationHttp "example/internal/input/application/http"
	inputApplicationHttpAdminAuthentication "example/internal/input/application/http/admin/authentication"
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
	ResourceClient *client.ResourceClient

	// HTTP server -Controller
	HttpAdminAuthenticationAuthenticator *inputApplicationHttpAdminAuthentication.AuthenticatorHandler

	// HTTP server -Middleware
	// Middleware 部分
	AdminAbstractMiddleware       *middlewareAdmin.AbstractMiddleware
	AdminAdminMiddleware          *middlewareAdmin.AdminMiddleware
	AdminAuthenticationMiddleware *middlewareAdmin.AuthenticationMiddleware
	AdminDecryptionMiddleware     *middlewareAdmin.DecryptionMiddleware
	AdminEncryptionMiddleware     *middlewareAdmin.EncryptionMiddleware
	AdminErrorMiddleware          *middlewareAdmin.ErrorMiddleware
	AdminLoggerMiddleware         *middlewareAdmin.LoggerMiddleware
	AdminNonexistentMiddleware    *middlewareAdmin.NonexistentMiddleware
	AdminRequestMiddleware        *middlewareAdmin.RequestMiddleware
	AdminResponseMiddleware       *middlewareAdmin.ResponseMiddleware
	AdminSignatureMiddleware      *middlewareAdmin.SignatureMiddleware
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
		usecaseApplicationAnyAdminAuthentication.NewAbstractUsecase,
		usecaseApplicationAnyAdminAuthentication.NewAuthenticatorUsecase,

		// client
		client.NewModel,
		client.NewResourceClient,

		// input-http
		inputApplicationHttp.NewAbstractHandler,

		inputApplicationHttpAdminAuthentication.NewAuthenticatorHandler,

		// Middleware 部分
		middlewareAdmin.NewAbstractMiddleware,
		middlewareAdmin.NewAdminMiddleware,
		middlewareAdmin.NewAuthenticationMiddleware,
		middlewareAdmin.NewDecryptionMiddleware,
		middlewareAdmin.NewEncryptionMiddleware,
		middlewareAdmin.NewErrorMiddleware,
		middlewareAdmin.NewLoggerMiddleware,
		middlewareAdmin.NewNonexistentMiddleware,
		middlewareAdmin.NewRequestMiddleware,
		middlewareAdmin.NewResponseMiddleware,
		middlewareAdmin.NewSignatureMiddleware,

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
	FacadeAbstract           *inputApplicationFacade.AbstractHandler
	FacadeTableScanner       *inputApplicationFacadeTable.ScannerHandler
	FacadeTableAuthenticator *inputApplicationFacadeRegister.AuthenticatorHandler

	// gRPC Facade Interceptor
	FacadeAdminErrorInterceptor          *interceptorFacadeGame.ErrorInterceptor
	FacadeAdminStatusInterceptor         *interceptorFacadeGame.StatusInterceptor
	FacadeAdminLoggerInterceptor         *interceptorFacadeGame.LoggerInterceptor
	FacadeAdminAuthenticationInterceptor *interceptorFacadeGame.AuthenticationInterceptor
}

func InitFacadeContainer() (*FacadeContainer, error) {
	wire.Build(

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewRsaHelper,

		// input-facade
		inputApplicationFacade.NewAbstractHandler,
		inputApplicationFacadeTable.NewScannerHandler,
		inputApplicationFacadeRegister.NewAuthenticatorHandler,

		// interceptor-facade
		interceptorFacadeGame.NewAbstractInterceptor,
		interceptorFacadeGame.NewErrorInterceptor,
		interceptorFacadeGame.NewStatusInterceptor,
		interceptorFacadeGame.NewLoggerInterceptor,
		interceptorFacadeGame.NewAuthenticationInterceptor,

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

	*usecaseApplicationAnyModel.AbstractUsecase
	usecasePortAnyModel.AdminUserUsecase

	// gRPC Resource server
	ResourceAbstract       *inputApplicationResource.AbstractHandler
	ResourceModelAdminUser *inputApplicationResourceModel.AdminUserHandler

	// gRPC Resource Interceptor
	ResourceAllInterceptor *interceptorResource.AllInterceptor
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
		outputApplicationMysql.NewAbstractRepository,
		outputApplicationMysqlModel.NewAdminUserRepository,

		// usecase
		usecaseApplicationAnyModel.NewAbstractUsecase,
		usecaseApplicationAnyModel.NewAdminUserUsecase,

		// input-resource
		inputApplicationResource.NewAbstractHandler,
		inputApplicationResourceModel.NewAdminUserHandler,

		// interceptor-resource
		interceptorResource.NewAbstractInterceptor,
		interceptorResource.NewAllInterceptor,

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
	*inputApplicationConsumer.AbstractHandler
	ConsumerAdminResourceAppUser *inputApplicationConsumerAdminResource.AppUserHandler
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
		outputApplicationMysql.NewAbstractRepository,
		outputApplicationMysqlModel.NewAppUserRepository,

		// usecase
		usecaseApplicationAnyAdminResource.NewAppUserUsecase,

		// input-consumer
		inputApplicationConsumer.NewAbstractHandler,
		inputApplicationConsumerAdminResource.NewAppUserHandler,

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
	*helper.JwtHelper

	// 排程 server
	CronAdminResourceAppUser             *inputApplicationCronAdminResource.AppUserHandler
	CronAdminAuthenticationAuthenticator *inputApplicationCronAdminAuthentication.AuthenticatorHandler
}

func InitCronContainer() (*CronContainer, error) {
	wire.Build(

		// bootstrap
		bootstrap.NewMysql,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewJwtHelper,

		// output
		outputApplicationMysql.NewAbstractRepository,
		outputApplicationMysqlModel.NewAppUserRepository,
		outputApplicationMysqlModel.NewAdminUserRepository,

		// usecase
		usecaseApplicationAnyAdminResource.NewAppUserUsecase,
		usecaseApplicationAnyAdminAuthentication.NewAbstractUsecase,
		usecaseApplicationAnyAdminAuthentication.NewAuthenticatorUsecase,

		// input-cron
		inputApplicationCron.NewAbstractHandler,
		inputApplicationCronAdminResource.NewAppUserHandler,
		inputApplicationCronAdminAuthentication.NewAuthenticatorHandler,

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
	*helper.JwtHelper

	// command
	*inputApplicationCommand.AbstractHandler
	CommandAdminReourceAppUser       *inputApplicationCommandAdminResource.AppUserHandler
	CommandAdminAuthenticationSignIn *inputApplicationCommandAdminAuthentication.AuthenticatorHandler
}

func InitCommandContainer() (*CommandContainer, error) {
	wire.Build(

		// bootstrap
		bootstrap.NewMysql,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewJwtHelper,

		// output
		outputApplicationMysql.NewAbstractRepository,
		outputApplicationMysqlModel.NewAppUserRepository,
		outputApplicationMysqlModel.NewAdminUserRepository,

		// usecase
		usecaseApplicationAnyAdminResource.NewAppUserUsecase,
		usecaseApplicationAnyAdminAuthentication.NewAbstractUsecase,
		usecaseApplicationAnyAdminAuthentication.NewAuthenticatorUsecase,

		// command
		inputApplicationCommand.NewAbstractHandler,
		inputApplicationCommandAdminResource.NewAppUserHandler,
		inputApplicationCommandAdminAuthentication.NewAuthenticatorHandler,

		wire.Struct(new(CommandContainer), "*"),
	)
	return nil, nil
}

type SourceContainer struct {

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper

	SourceAnnouncementLottery *inputApplicationSourceAnnouncement.LotteryHandler
}

func InitSourceContainer() (*SourceContainer, error) {
	wire.Build(

		// bootstrap

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,

		outputApplicationMemory.NewAbstractRepository,
		outputApplicationMemoryModel.NewLotteryRepository,
		usecaseApplicationAnyAnnouncement.NewAbstractUsecase,
		usecaseApplicationAnyAnnouncement.NewLotteryUsecase,

		inputApplicationSource.NewAbstractHandler,
		inputApplicationSourceAnnouncement.NewLotteryHandler,

		wire.Struct(new(SourceContainer), "*"),
	)
	return nil, nil
}

type DaemonContainer struct {

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper

	// Clients
	SourceClient *client.SourceClient

	DaemonWatcherSourceAnnouncementLottery *inputApplicationDaemonWatcherSource.AnnouncementLotteryHandler
}

func InitDaemonContainer() (*DaemonContainer, error) {
	wire.Build(

		// bootstrap
		bootstrap.NewSource,
		bootstrap.NewRedis,

		// client
		client.NewAnnouncement,
		client.NewSourceClient,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewCacheHelper,

		// output
		outputApplicationCache.NewAbstractRepository,
		outputApplicationCacheModel.NewLotteryRepository,

		// usecase
		usecaseApplicationAnyWatcherSource.NewAbstractUsecase,
		usecaseApplicationAnyWatcherSource.NewAnnouncementLotteryUsecase,

		// input-daemon
		inputApplicationDaemon.NewAbstractHandler,
		inputApplicationDaemonWatcherSource.NewAnnouncementLotteryHandler,

		wire.Struct(new(DaemonContainer), "*"),
	)
	return nil, nil
}
