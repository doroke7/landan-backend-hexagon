package announcement

import (
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	inputApplicationSource "example/internal/input/application/source"
	pbSourceAnnouncement "example/pb/source/announcement"
)

type LotteryHandler struct {
	pbSourceAnnouncement.UnimplementedLotteryServer
	*inputApplicationSource.AbstractHandler
}

func NewLotteryHandler(oAbstractHandler *inputApplicationSource.AbstractHandler) *LotteryHandler {
	return &LotteryHandler{
		AbstractHandler: oAbstractHandler,
	}
}

func (oSelf *LotteryHandler) Watch(oReq *pbSourceAnnouncement.LotteryWatchRequest, oStream grpc.ServerStreamingServer[pbSourceAnnouncement.LotteryWatchReply]) error {

	oTicker := time.NewTicker(3 * time.Second)
	defer oTicker.Stop()

	var iCount int32 = 1

	for {
		select {
		// 情境 A：Client 斷線或 Cancel 請求
		case <-oStream.Context().Done():
			return oStream.Context().Err()

		// 情境 B：時間到，準備發送資料
		case oT := <-oTicker.C:
			iCount++
			// 模擬產生開獎資料 (實務上這裡會是呼叫 UseCase / DB 撈資料)
			oResponse := &pbSourceAnnouncement.LotteryWatchReply{
				Id:      iCount,
				Round:   fmt.Sprintf("20260721-%03d", iCount),
				Time:    oT.Unix(),
				Numbers: []int32{10, 20, 30},
			}

			// 透過 stream.Send() 推送資料給 Client
			if err := oStream.Send(oResponse); err != nil {
				log.Printf("資料推送失敗: %v", err)
				return err // 發送失敗（通常是網路斷開），中斷 return
			}

			log.Printf("成功推送期號: %s", oResponse.Round)
		}
	}
}
