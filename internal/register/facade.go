package register

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	pkg "example/pkg"

	container "example/container"

	pbFacadeAdminAuthentication "example/pb/facade/admin/authentication"
	pbFacadeRegister "example/pb/facade/register"
	pbFacadeTable "example/pb/facade/table"
)

func facadeInterceptors(oContainer *container.FacadeContainer) grpc.UnaryServerInterceptor {

	aGameInterceptors := []grpc.UnaryServerInterceptor{
		oContainer.FacadeGameErrorInterceptor.Handle(),
		oContainer.FacadeGameStatusInterceptor.Handle(),
		oContainer.FacadeGameLoggerInterceptor.Handle(),
	}

	aAdminInterceptors := []grpc.UnaryServerInterceptor{
		oContainer.FacadeAdminErrorInterceptor.Handle(),
		oContainer.FacadeAdminStatusInterceptor.Handle(),
		oContainer.FacadeAdminLoggerInterceptor.Handle(),
		oContainer.FacadeAdminSignatureInterceptor.Handle(),
		oContainer.FacadeAdminDecryptionInterceptor.Handle(),
	}

	oRouter := pkg.NewGrpcRouter()

	oRouter.Group("/pb.facade.game.", aGameInterceptors...)
	oRouter.Group("/pb.facade.game.authentication.")
	oRouter.Group("/pb.facade.game.resource.", oContainer.FacadeGameAuthenticationInterceptor.Handle())

	oRouter.Group("/pb.facade.admin.", aAdminInterceptors...)

	return oRouter.Build()
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
	pbFacadeAdminAuthentication.RegisterAuthenticatorServer(oGrpcServer, oContainer.FacadeAdminAuthenticationAuthenticator)

	return oGrpcServer
}
