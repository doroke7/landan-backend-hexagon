package handler

import (
	"net/http"
	"strconv"

	HttpAdmin "example/internal/input/http/admin"
	"example/internal/usecase/port"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	*HttpAdmin.AbstractHandler
	userUsecase port.UserUsecase // 不能把把 userUsecase（port.UserUsecase）塞進 AbstractHandler 確實不對——UserUsecase 是 driving port，是每個 handler 各自要注入的業務依賴，不是「抽象共用的技術基礎設施」
}

func NewUserHandler(oUserUsecase port.UserUsecase, oAbstractHandler *HttpAdmin.AbstractHandler) *UserHandler {
	return &UserHandler{
		AbstractHandler: oAbstractHandler,
		userUsecase:     oUserUsecase,
	}
}

func (oSelf *UserHandler) AddUser(oContext *gin.Context) {

	name := oContext.Query("name")

	user, err := oSelf.userUsecase.AddUserByName(name)
	if err != nil {
		oContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	oContext.JSON(http.StatusOK, user)
}

func (oSelf *UserHandler) ShowUser(oContext *gin.Context) {

	idStr := oContext.Query("id")

	id, _ := strconv.Atoi(idStr)

	user, err := oSelf.userUsecase.ShowUserById(id)
	if err != nil {
		oContext.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	oContext.JSON(http.StatusOK, user)
}
