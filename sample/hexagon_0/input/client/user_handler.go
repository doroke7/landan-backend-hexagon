package client

import (
	"context"
	"io"
	"log"

	"google.golang.org/grpc"

	"example/input/abstract"
	"example/input/port"
	pb "example/pb/client"
)

// UserHandler 主動訂閱外部 gRPC stream server 推送的 User 事件，
// 每收到一筆就呼叫 usecase 新增用戶。
type UserHandler struct {
	*abstract.AbstractHandler
	userUsecase port.UserUsecase
	client      pb.UserServiceClient
}

func NewUserHandler(conn *grpc.ClientConn, useCase port.UserUsecase, oAbstractHandler *abstract.AbstractHandler) *UserHandler {
	return &UserHandler{
		AbstractHandler: oAbstractHandler,
		userUsecase:     useCase,
		client:          pb.NewUserServiceClient(conn),
	}
}

func (oSelf *UserHandler) Start(ctx context.Context) error {

	stream, err := oSelf.client.SubscribeUsers(ctx, &pb.SubscribeUsersRequest{})
	if err != nil {
		return err
	}

	for {
		user, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if _, err := oSelf.userUsecase.CreateUser(user.GetName()); err != nil {
			log.Printf("create user failed: %v", err)
		}
	}
}
