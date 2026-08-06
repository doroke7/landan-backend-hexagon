package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbResourceModel "example/pb/resource/model"

	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
)

type AdminUserHandler struct {
	pbResourceModel.UnimplementedAdminUserServer
	*inputApplicationResource.AbstractHandler
	usecasePortAnyModel.AdminUserUsecase
}

func NewAdminUserHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oAdminUserUsecase usecasePortAnyModel.AdminUserUsecase) *AdminUserHandler {
	return &AdminUserHandler{
		AbstractHandler:  oAbstractHandler,
		AdminUserUsecase: oAdminUserUsecase,
	}
}

func (oSelf *AdminUserHandler) ShowOneByName(oContext context.Context, oReq *pbResourceModel.OneAdminUserRequest) (*pbResourceModel.OneAdminUerResponse, error) {

	oAdminUser, err := oSelf.AdminUserUsecase.ShowOneByName(oReq.Name)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pbResourceModel.OneAdminUerResponse{
		Id:       int32(oAdminUser.Id),
		Name:     oAdminUser.Name,
		Password: oAdminUser.Password,
	}, nil

}
