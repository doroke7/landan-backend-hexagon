package register

import (
	"google.golang.org/grpc"

	container "example/internal/container"

	pbFacadeGame "example/pb/facade/game"
	pbFacadeTable "example/pb/facade/table"
)

func FacadeInit(oContainer *container.Container) *grpc.Server {
	oGrpcServer := grpc.NewServer()

	pbFacadeTable.RegisterScannerServer(oGrpcServer, oContainer.FacadeTableScannerUser)
	pbFacadeGame.RegisterUserServiceServer(oGrpcServer, oContainer.FacadeGameUser)

	return oGrpcServer
}
