package command

import (
	"github.com/spf13/cobra"

	"example/internal/usecase/port"
	"log"
)

func NewUserHandler(oUserUsecase port.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: oUserUsecase,
	}
}

type UserHandler struct {
	userUsecase port.UserUsecase // UserUsecase 是 driving port，是每個 handler 各自要注入的業務依賴，不是「抽象共用的技術基礎設施」，不能塞進 AbstractHandler

}

func (oSelf *UserHandler) AddUser() *cobra.Command {

	var oUserAddUserCommand = &cobra.Command{
		Use:   "User-AddUser",
		Short: "User-AddUser 相關命令",
		Run: func(oCmd *cobra.Command, args []string) {
			sName := "Tom"

			if _, err := oSelf.userUsecase.AddUserByName(sName); err != nil {
				log.Printf("create user failed: %v", err)
				return
			}
		},
	}
	return oUserAddUserCommand

}
