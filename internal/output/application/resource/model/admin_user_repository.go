package resource

import (
	"context"
	"errors"

	client "example/internal/client"
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResourceModel "example/pb/resource/model"

	outputApplicationResource "example/internal/output/application/resource"
)

type AdminUserRepository struct {
	AdminUserModelClient pbResourceModel.AdminUserClient
	*outputApplicationResource.AbstractRepository
}

func NewAdminUserRepository(oResourceClient *client.ResourceClient, oAbstractRepository *outputApplicationResource.AbstractRepository) outputPortAnyModel.AdminUserRepository {
	return &AdminUserRepository{
		AdminUserModelClient: oResourceClient.Model.AdminUser,
		AbstractRepository:   oAbstractRepository,
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

// 一點技巧，可以先 回傳 nil, Error 來先不寫實作，不然一次全部 output 補齊太麻煩了
// ShowOneById 目前 Resource gRPC service 沒有對應的 RPC，先不支援。
func (oSelf *AdminUserRepository) ShowOneById(id uint) (*domain.AdminUser, error) {
	return nil, errors.New("not supported by resource")
}
