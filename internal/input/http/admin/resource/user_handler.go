package handler

import (
	"net/http"
	"strconv"

	HttpAdmin "example/internal/input/http/admin"
	"example/internal/usecase/port/model"
	pkg "example/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		pkg.Logger(pkg.HttpAdmin).Error("AddUser 失敗",
			zap.String("name", name),
			zap.Error(err),
		)
		oContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pkg.Logger(pkg.HttpAdmin).Info("AddUser 成功",
		zap.Int("id", user.ID),
		zap.String("name", user.Name),
	)
	oContext.JSON(http.StatusOK, user)
}

func (oSelf *UserHandler) ShowUser(oContext *gin.Context) {

	idStr := oContext.Query("id")

	id, _ := strconv.Atoi(idStr)

	user, err := oSelf.userUsecase.ShowUserById(id)
	if err != nil {
		pkg.Logger(pkg.HttpAdmin).Warn("ShowUser 失敗",
			zap.Int("id", id),
			zap.Error(err),
		)
		oContext.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	pkg.Logger(pkg.HttpAdmin).Info("ShowUser 成功",
		zap.Int("id", user.ID),
		zap.String("name", user.Name),
	)
	oContext.JSON(http.StatusOK, user)
}
