package service

import (
	"context"

	"example/internal/input/port"
	resource "example/internal/input/resource"
	pb "example/pb/resource/model"
)

type AdminUserHandler struct {
	pb.UnimplementedAdminUserServer
	*resource.AbstractHandler
	AdminUserUsecase port.AdminUserUsecase
}

func NewAdminUserHandler(oAbstractHandler *resource.AbstractHandler) *AdminUserHandler {
	return &AdminUserHandler{
		AbstractHandler: oAbstractHandler,
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
