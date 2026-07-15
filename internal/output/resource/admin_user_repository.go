package resource

import (
	"context"

	client "example/internal/client"
	"example/internal/domain"
	"example/internal/output/port"
	pbResourceModel "example/pb/resource/model"
)

type AdminUserRepository struct {
	AdminUserClient pbResourceModel.AdminUserClient
}

func NewAdminUserRepository(oResourceClient *client.ResourceClient) port.AdminUserRepository {
	return &AdminUserRepository{
		AdminUserClient: oResourceClient.Model.AdminUser,
	}
}

func (oSelf *AdminUserRepository) ShowOneByName(sName string) (*domain.AdminUser, error) {
	oResp, err := oSelf.AdminUserClient.ShowOneByName(
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
