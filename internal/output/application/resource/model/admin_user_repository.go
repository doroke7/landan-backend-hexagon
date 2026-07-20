package resource

import (
	"context"

	client "example/internal/client"
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResourceModel "example/pb/resource/model"
)

type AdminUserRepository struct {
	AdminUserModelClient pbResourceModel.AdminUserClient
}

func NewAdminUserRepository(oResourceClient *client.ResourceClient) outputPortAnyModel.AdminUserRepository {
	return &AdminUserRepository{
		AdminUserModelClient: oResourceClient.Model.AdminUser,
	}
}

func (oSelf *AdminUserRepository) ShowOneByName(sName string) (*domain.AdminUser, error) {

	oResp, err := oSelf.AdminUserModelClient.ShowOneByName(
		context.Background(),
		&pbResourceModel.OneAdminUserRequest{Name: sName},
	)

	if err != nil {
		return nil, err
	}

	return &domain.AdminUser{
		Id:       uint(oResp.GetId()),
		Name:     oResp.GetName(),
		Password: oResp.GetPassword(),
	}, nil
}
