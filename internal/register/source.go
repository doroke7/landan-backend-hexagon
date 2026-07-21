package register

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	pbSourceAnnouncement "example/pb/source/announcement"

	container "example/internal/container"
	pkg "example/pkg"
)

func sourceInterceptors(_ *container.SourceContainer) grpc.UnaryServerInterceptor {

	// aBase 全局攔截器：目前 resource 只需要 facade -> resource 的 Basic Auth 驗證
	aBase := []grpc.UnaryServerInterceptor{
		// oContainer.ResourceAllInterceptor.Handle(),
	}

	return pkg.NewRouter(aBase...).
		// database 群組目前只需要全局驗證；未來個別 resource 服務需要額外攔截器
		// （例如特定表要多一層權限檢查）時，在這裡加一個更具體的 Group 前綴即可
		Build()
}

func SourceInit(oContainer *container.SourceContainer) *grpc.Server {

	oGrpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(sourceInterceptors(oContainer)),
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
	pbSourceAnnouncement.RegisterLotteryServer(oGrpcServer, oContainer.SourceAnnouncementLottery)

	return oGrpcServer
}
