package facade

import (
	"context"

	inputFacade "example/internal/input/application/facade"
	"example/internal/usecase/port/facade/model"
	pb "example/pb/facade/game"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	*inputFacade.AbstractHandler
	userUsecase port.UserUsecase // 不能把把 userUsecase（port.UserUsecase）塞進 AbstractHandler 確實不對——UserUsecase 是 driving port，是每個 handler 各自要注入的業務依賴，不是「抽象共用的技術基礎設施」
}

func NewUserHandler(useCase port.UserUsecase, oAbstractHandler *inputFacade.AbstractHandler) *UserHandler {
	return &UserHandler{
		AbstractHandler: oAbstractHandler,
		userUsecase:     useCase,
	}
}

func (oSelf *UserHandler) AddUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	user, err := oSelf.userUsecase.AddUserByName(req.GetName())
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:   int64(user.ID),
			Name: user.Name,
		},
	}, nil
}

func (oSelf *UserHandler) ShowUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {

	user, err := oSelf.userUsecase.ShowUserById(int(req.GetId()))
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:   int64(user.ID),
			Name: user.Name,
		},
	}, nil
}
