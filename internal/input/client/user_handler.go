package client

import (
	"context"
	"io"

	"example/internal/usecase/port/model"
	pb "example/pb/client"
	pkg "example/pkg"

	"go.uber.org/zap"
)

// UserHandler 主動訂閱外部 gRPC stream server 推送的 User 事件，
// 每收到一筆就呼叫 usecase 新增用戶。
type UserHandler struct {
	*AbstractHandler
	userUsecase port.UserUsecase // 不能把把 userUsecase（port.UserUsecase）塞進 AbstractHandler 確實不對——UserUsecase 是 driving port，是每個 handler 各自要注入的業務依賴，不是「抽象共用的技術基礎設施」
}

func NewUserHandler(useCase port.UserUsecase, oAbstractHandler *AbstractHandler) *UserHandler {
	return &UserHandler{
		AbstractHandler: oAbstractHandler,
		userUsecase:     useCase,
	}
}

func (oSelf *UserHandler) AddUser(ctx context.Context) error {

	stream, err := oSelf.Client.User.SubscribeUsers(ctx, &pb.SubscribeUsersRequest{})
	if err != nil {
		pkg.Logger(pkg.Client).Error("SubscribeUsers 失敗",
			zap.Error(err),
		)
		return err
	}

	for {
		// 可以 把 channel 丟到一個
		user, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			pkg.Logger(pkg.Client).Error("stream.Recv 失敗",
				zap.Error(err),
			)
			return err
		}

		if _, err := oSelf.userUsecase.AddUserByName(user.GetName()); err != nil {
			pkg.Logger(pkg.Client).Error("AddUser 失敗",
				zap.String("name", user.GetName()),
				zap.Error(err),
			)
			continue
		}

		pkg.Logger(pkg.Client).Info("AddUser 成功",
			zap.String("name", user.GetName()),
		)
	}
}
