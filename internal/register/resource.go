package register

import (
	"google.golang.org/grpc"

	container "example/internal/container"

	pbResourceModel "example/pb/resource/model"
)

func ResourceInit(oContainer *container.ResourceContainer) *grpc.Server {
	oGrpcServer := grpc.NewServer()

	pbResourceModel.RegisterAdminUserServer(oGrpcServer, oContainer.)

	return oGrpcServer
}
