package register

import (
	"google.golang.org/grpc"

	container "example/internal/container"

	pbResourceModel "example/pb/resource/model"
)

func ResourceInit(oContainer *container.ResourceContainer, oGrpcServer *grpc.Server) *grpc.Server {

	pbResourceModel.RegisterAdminUserServer(oGrpcServer, oContainer.ResourceModelAdminUser)

	return oGrpcServer
}
