package register

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	pbResourceModel "example/pb/resource/model"

	container "example/container"
	pkg "example/pkg"
)

func resourceInterceptors(oContainer *container.ResourceContainer) grpc.UnaryServerInterceptor {

	// aBase 全局攔截器：目前 resource 只需要 facade -> resource 的 Basic Auth 驗證
	aBase := []grpc.UnaryServerInterceptor{
		oContainer.ResourceAllInterceptor.Handle(),
	}

	return pkg.NewRouter(aBase...).
		// database 群組目前只需要全局驗證；未來個別 resource 服務需要額外攔截器
		// （例如特定表要多一層權限檢查）時，在這裡加一個更具體的 Group 前綴即可
		Group("/resource.").
		Build()
}

func ResourceInit(oContainer *container.ResourceContainer) *grpc.Server {

	oGrpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(resourceInterceptors(oContainer)),
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

	pbResourceModel.RegisterAdminUserServer(oGrpcServer, oContainer.ResourceModelAdminUser)

	return oGrpcServer
}
