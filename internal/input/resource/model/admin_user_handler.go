package service

import (
	"context"

	inputResource "example/internal/input/resource"
	port "example/internal/usecase/port/resource/model"
	pb "example/pb/resource/model"
)

type AdminUserHandler struct {
	pb.UnimplementedAdminUserServer
	*inputResource.AbstractHandler
	port.AdminUserUsecase
}

func NewAdminUserHandler(oAbstractHandler *inputResource.AbstractHandler, oAdminUserUsecase port.AdminUserUsecase) *AdminUserHandler {
	return &AdminUserHandler{
		AbstractHandler:  oAbstractHandler,
		AdminUserUsecase: oAdminUserUsecase,
	}
}

func (oSelf *AdminUserHandler) ShowOneByName(oContext context.Context, oReq *pb.OneAdminUserRequest) (*pb.OneAdminUerResponse, error) {

	oAdminUser, err := oSelf.AdminUserUsecase.ShowOneByName(oReq.Name)
	if err != nil {
		return nil, err
	}

	return &pb.OneAdminUerResponse{
		Id:       int32(oAdminUser.Id),
		Name:     oAdminUser.Name,
		Password: oAdminUser.Password,
	}, nil

}
