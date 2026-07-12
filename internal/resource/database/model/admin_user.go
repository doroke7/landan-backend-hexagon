package service

import (
	"context"

	model "landan-backend-grpc/internal/model/database"
	pb "landan-backend-grpc/pb/resource/database/model"
	"landan-backend-grpc/pkg"
)

type AdminUserService struct {
	pb.UnimplementedAdminUserServer
	*AbstractService
}

func NewAdminUserService(oAbstractService *AbstractService) *AdminUserService {
	return &AdminUserService{
		AbstractService: oAbstractService,
	}
}

func (oSelf *AdminUserService) ShowOneByName(oContext context.Context, oReq *pb.AdminUserShowOneByNameInput) (*pb.AdminUserShowOneByNameResult, error) {

	oAdminUserModel, err := pkg.DatabaseInit[*model.AdminUserModel](oSelf.DatabaseFactory, "AdminUser", "1")
	if err != nil {
		return nil, err
	}

	oAdminUser, oErr := oAdminUserModel.ShowOneByName(oReq.Name)

	if nil != oErr {
		return nil, oErr
	}

	return &pb.AdminUserShowOneByNameResult{
		One: &pb.AdminUserOne{
			Id:       oAdminUser.Id,
			Name:     oAdminUser.Name,
			Password: oAdminUser.Password,
		},
	}, nil

}
