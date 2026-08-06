package grpc

import (
	"context"

	"example/input/abstract"
	"example/input/port"
	pb "example/pb/grpc"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	*abstract.AbstractHandler
	userUsecase port.UserUsecase
}

func NewUserHandler(useCase port.UserUsecase, oAbstractHandler *abstract.AbstractHandler) *UserHandler {
	return &UserHandler{
		AbstractHandler: oAbstractHandler,
		userUsecase:     useCase,
	}
}

func (oSelf *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	user, err := oSelf.userUsecase.CreateUser(req.GetName())
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

func (oSelf *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {

	user, err := oSelf.userUsecase.GetUser(int(req.GetId()))
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
