package service

import (
	"context"

	resource "example/internal/input/resource"
	"example/internal/usecase/port"
	pb "example/pb/resource/model"
)

type AdminUserHandler struct {
	pb.UnimplementedAdminUserServer
	*resource.AbstractHandler
	port.AdminUserUsecase
}

func NewAdminUserHandler(oAbstractHandler *resource.AbstractHandler, oAdminUserUsecase port.AdminUserUsecase) *AdminUserHandler {
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
