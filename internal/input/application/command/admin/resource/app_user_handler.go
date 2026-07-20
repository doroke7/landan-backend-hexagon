package command

import (
	"log"

	"github.com/spf13/cobra"

	inputApplicationCommand "example/internal/input/application/command"
	port "example/internal/usecase/port/any/admin/resource"
)

type AppUserHandler struct {
	*inputApplicationCommand.AbstractHandler
	appUserUsecase port.AppUserUsecase
}

func NewAppUserHandler(oAppUserUsecase port.AppUserUsecase, oAbstractHandler *inputApplicationCommand.AbstractHandler) *AppUserHandler {
	return &AppUserHandler{
		AbstractHandler: oAbstractHandler,
		appUserUsecase:  oAppUserUsecase,
	}
}

func (oSelf *AppUserHandler) IncreaseBalance() *cobra.Command {

	var iId uint
	var iAmount uint

	var oAppUserIncreaseBalanceCommand = &cobra.Command{
		Use:   "Admin-Resource-AppUser-IncreaseBalance",
		Short: "AppUser-IncreaseBalance 相關命令",
		Run: func(oCmd *cobra.Command, args []string) {
			if _, err := oSelf.appUserUsecase.IncreaseBalance(iId, iAmount); err != nil {
				log.Printf("increase balance failed: %v", err)
				return
			}
		},
	}

	oAppUserIncreaseBalanceCommand.Flags().UintVar(&iId, "id", 1, "AppUser 的 id")
	oAppUserIncreaseBalanceCommand.Flags().UintVar(&iAmount, "amount", 10, "要增加的餘額")

	return oAppUserIncreaseBalanceCommand
}
