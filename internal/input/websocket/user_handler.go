package websocket

import (
	"net/http"

	"example/internal/usecase/resource/port/model"
	pkg "example/pkg"

	"go.uber.org/zap"
)

type createUserMessage struct {
	Name string `json:"name"`
}

type UserHandler struct {
	*AbstractHandler
	userUsecase port.UserUsecase // UserUsecase 是 driving port，是每個 handler 各自要注入的業務依賴，不是「抽象共用的技術基礎設施」，不能塞進 AbstractHandler
}

func NewUserHandler(oUserUsecase port.UserUsecase, oAbstractHandler *AbstractHandler) *UserHandler {
	return &UserHandler{
		AbstractHandler: oAbstractHandler,
		userUsecase:     oUserUsecase,
	}
}

// AddUser 職責跟其他 adapter 的 AddUser 一樣單純：把連線升級成 websocket 之後，
// 每收到一筆訊息就轉呼叫 usecase 新增用戶，再把結果寫回去。
func (oSelf *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	oConn, err := oSelf.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		pkg.Logger(pkg.Websocket).Error("websocket upgrade 失敗",
			zap.Error(err),
		)
		return
	}
	defer oConn.Close()

	for {
		var payload createUserMessage
		if err := oConn.ReadJSON(&payload); err != nil {
			pkg.Logger(pkg.Websocket).Error("ReadJSON 失敗",
				zap.Error(err),
			)
			return
		}

		user, err := oSelf.userUsecase.AddUserByName(payload.Name)
		if err != nil {
			pkg.Logger(pkg.Websocket).Error("AddUser 失敗",
				zap.String("name", payload.Name),
				zap.Error(err),
			)
			continue
		}

		pkg.Logger(pkg.Websocket).Info("AddUser 成功",
			zap.String("name", payload.Name),
		)

		if err := oConn.WriteJSON(user); err != nil {
			pkg.Logger(pkg.Websocket).Error("WriteJSON 失敗",
				zap.Error(err),
			)
			return
		}
	}
}
