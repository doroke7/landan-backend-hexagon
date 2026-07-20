package register

import (
	pkg "example/pkg"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	container "example/internal/container"

	pbFacadeRegister "example/pb/facade/register"
	pbFacadeTable "example/pb/facade/table"
)

func facadeInterceptors(oContainer *container.FacadeContainer) grpc.UnaryServerInterceptor {

	// aBase 全局攔截器：任何 group 都會套用，等同 Hyperf 的全局 middleware
	aBase := []grpc.UnaryServerInterceptor{
		oContainer.FacadeAdminErrorInterceptor.Handle(),
		oContainer.FacadeAdminStatusInterceptor.Handle(),
		oContainer.FacadeAdminLoggerInterceptor.Handle(),
	}

	return pkg.NewRouter(aBase...).
		// 登入接口：不需要額外驗證，只套用全局攔截器
		Group("/facae.admin.authentication.").
		// 資源接口：全局攔截器 + Token 驗證（proto 尚未生成，先預留 group，
		// 之後 admin/resource 的 service 上線會自動吃到這條規則）
		Group("/facae.admin.resource.", oContainer.FacadeAdminAuthenticationInterceptor.Handle()).
		Build()
}

func FacadeInit(oContainer *container.FacadeContainer) *grpc.Server {

	oGrpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(facadeInterceptors(oContainer)),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    1 * time.Second,
				Timeout: 5 * time.Second,
			},
		),
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             10 * time.Second,
				PermitWithoutStream: true,
			},
		),
	)
	pbFacadeTable.RegisterScannerServer(oGrpcServer, oContainer.FacadeTableScanner)
	pbFacadeRegister.RegisterAuthenticatorServer(oGrpcServer, oContainer.FacadeTableAuthenticator)

	return oGrpcServer
}
