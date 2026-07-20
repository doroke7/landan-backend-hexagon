package cron

type AppUserHandler struct {
	*AbstractHandler
}

func NewAppUserHandler(oAbstractHandler *AbstractHandler) *AppUserHandler {
	return &AppUserHandler{
		AbstractHandler: oAbstractHandler,
	}
}

func (oSelf *AppUserHandler) IncreaseBalance() {

}
