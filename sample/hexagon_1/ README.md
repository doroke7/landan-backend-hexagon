## 六角架構圖

+----------------------------------------------------------------+
|                            input                               |
|        HTTP / gRPC / CLI / Cron / WebSocket / GraphQL          |
+----------------------------------------------------------------+
                               |
                               v
                    +----------------------+
                    |     input_port       |
                    +----------------------+
                               ^
                               |
                    +----------------------+      +----------------------+
                    |                      |----->|                      |
                    |                      |      |                      |
                    |                      |      |                      |
                    |                      |      |        domain        |
                    |       use_case       |      |                      |
                    |                      |      |                      |
                    |                      |      |                      |
                    |                      |      |                      |
                    +----------+-----------+      +----------------------+
                               |
                               v
                    +----------------------+
                    |     output_port      |
                    +----------------------+
                               ^
                               |
+----------------------------------------------------------------+
|                            output                              |
|      MySQL / Redis / Kafka / S3 / MQ / Third-party API         |
+----------------------------------------------------------------+

## 六角架構核心優點
1. 可以同時輸入 http grpc cron command 但是共用一個業務邏輯 usecase

## 目錄結構

```
.
├── main.go                    # 進入點，實際邏輯委派給 cmd.Execute()
├── cmd/                       # cobra 指令，每個檔案對應一個可獨立啟動的服務
│   ├── root.go                #   root command，Execute() 供 main.go 呼叫
│   ├── http.go                #   啟動 HTTP 服務
│   ├── grpc.go                #   啟動 gRPC 服務
│   ├── consumer.go            #   啟動 AMQP consumer
│   ├── client.go               #   啟動 gRPC client stream 訂閱
│   ├── cron.go                #   啟動排程服務
│   └── websocket.go           #   啟動 websocket 服務
│
├── internal/
│   ├── bootstrap/             # 讀 CONFIG、建立各種基礎設施連線（mysql / redis / amqp / mongo / grpc client）
│   ├── domain/                # 領域物件（entity），跟任何框架、資料庫無關
│   ├── helper/                # 通用工具（AES、RSA、快取讀寫……），跟業務邏輯無關可到處注入
│   │
│   ├── input/                 # Input Adapter：把外部請求轉換成呼叫 usecase
│   │   ├── port/               #   input port，usecase 對外暴露的介面（driving port）
│   │   ├── http/               #   HTTP handler
│   │   ├── grpc/                #   gRPC server handler
│   │   ├── client/              #   gRPC client（訂閱外部 stream）
│   │   ├── consumer/            #   AMQP consumer handler
│   │   ├── cron/                #   排程任務 handler
│   │   ├── websocket/           #   websocket handler
│   │   └── command/             #   CLI 指令 handler
│   │   （每個 adapter 底下都有自己獨立的 abstract_handler.go，
│   │    彼此不共用，只共用 usecase 這個核心業務邏輯）
│   │
│   ├── usecase/                # 業務邏輯本體，只依賴 input/port、output/port，不依賴任何 adapter
│   │
│   ├── output/                 # Output Adapter：usecase 依賴的下游資源實作
│   │   ├── port/                #   output port，usecase 依賴的介面（driven port）
│   │   ├── mysql/                #   MySQL 實作（gorm）
│   │   ├── cache/                #   裝飾器（Decorator），包住 mysql 實作，加上 redis 讀寫快取
│   │   ├── memory/               #   記憶體實作（測試/範例用）
│   │   └── producer/             #   訊息生產者實作
│   │
│   ├── register/               # 組裝層：把 container 生好的 handler 註冊到對應的 server/router 上
│   │                            #   （http.HandleFunc / grpc.RegisterXxxServer / cron.AddFunc ...），
│   │                            #   cmd/ 只管呼叫 XxxInit 拿到 server 物件再 Serve，不碰組裝細節
│   │
│   └── container/              # wire 組裝根：wire.go 手寫、wire_gen.go 自動產生，別手改後者
│
├── pkg/                        # 跟 domain 無關、可重用的通用元件
│   ├── consumer_router.go       #   queue name -> handler 的路由表（AMQP 沒有內建路由機制）
│   ├── client_router.go         #   多個 client-side 訂閱方法的並行啟動器
│   ├── websocket_router.go      #   websocket 路由的路徑前綴分組（模仿 gin Group）
│   └── aop.go                   #   泛型 Cacheable / CachePut / CacheEvict，AOP 風格的快取包裝
│
├── config/                     # viper 讀取的 yaml 設定檔，一個檔案對應一個頂層命名空間
└── pb/                          # protoc 產生的程式碼，對應 proto/ 底下的定義
```

依賴方向永遠是「外層指向內層」：`input adapter → input/port → usecase → output/port ← output adapter`，
`usecase` 完全不知道自己被 http 還是 grpc 還是 cron 呼叫，也不知道資料到底存在 mysql 還是 redis。