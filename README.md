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
1. 可以同時輸入 http / grpc / cron / consumer / websocket / client-stream / command，但共用同一套 usecase 業務邏輯。
2. 每種輸入各自只組裝自己需要的依賴（見下方「每個服務獨立 container」），不會因為要跑 `cron` 就順便把 gRPC client、AMQP 都連上。

## 目錄結構

```
.
├── main.go                        # 進入點，只呼叫 cmd.Execute()
├── cmd/                           # cobra 指令，每個檔案對應一個可獨立啟動的服務／進程
│   ├── root.go                    #   root command，Execute() 供 main.go 呼叫
│   ├── facade.go                  #   啟動 facade gRPC 服務（對外入口）
│   ├── resource.go                #   啟動 resource gRPC 服務（資料服務，僅供 facade / http 呼叫）
│   ├── http.go                    #   啟動 HTTP（Gin）服務
│   ├── consumer.go                #   啟動 AMQP consumer
│   ├── client.go / socket.go      #   啟動 gRPC client-side stream 訂閱
│   ├── cron.go                    #   啟動排程服務
│   ├── websocket.go               #   啟動 websocket 服務
│   └── command.go                 #   啟動一次性 CLI 指令
│
├── internal/
│   ├── bootstrap/                 # 讀 CONFIG、建立各種基礎設施連線（mysql / redis / amqp / mongo / grpc client）
│   ├── domain/                    # 領域物件（entity）：AdminUser、AppUser、User，跟框架、資料庫無關
│   ├── helper/                    # 通用工具（AES、RSA、JWT、Cache 讀寫……），跟業務邏輯無關可到處注入
│   ├── client/                    # 對外部 gRPC stream server / resource 服務的 client 封裝
│   │
│   ├── input/                     #   協議輸入端（driving adapter），只有實作，沒有介面
│   │   └── application/           #   不同協議輸入端，不同的端且相同的相對目錄 代表同一個業務輸入。
│   │       ├── facade/            #   對外 gRPC 入口
│   │       ├── resource/          #   resource 內部 gRPC 服務（僅供 facade / http 呼叫）
│   │       ├── http/              #   對外 Http 入口
│   │       ├── client/            #   gRPC client（訂閱外部 stream）
│   │       ├── consumer/          #   AMQP consumer
│   │       ├── cron/              #   排程任務
│   │       ├── websocket/         #   websocket 入口
│   │       └── command/           #   CLI 指令
│   │       （每個 adapter 底下都有自己獨立的 abstract_handler.go，彼此不共用；
│   │        adapter 內部依 leaf 功能再分 admin/resource、admin/authentication 這種子資料夾——
│   │        這是跨層的對應 key，不是字面語意：input/<adapter>/admin/resource、
│   │        usecase/application/any/admin/resource、usecase/port/any/admin/resource
│   │        三者相對路徑相同，代表同一條 usecase 邏輯，不代表跟 HTTP 後台路由有關）
│   │
│   ├── middleware/
│   │   └── admin/                    # HTTP 專用 middleware 鏈
│   │
│   ├── interceptor/                  # gRPC 專用攔截器鏈
│   │   ├── facade/
│   │   │   └── game/
│   │   └── resource/
│   │
│   ├── usecase/                   # 商務案例：實作 + 端口介面
│   │   ├── application/
│   │   │   └── any/                       #  「any」表示這份 usecase 不綁定特定 adapter，可以被多個 driving
│   │   │       │                          # usecase端一種命名結構對應不同input 但是多種輸入
│   │   │       ├── admin/                 #  後台服務
│   │   │       │   ├── authentication/    #  登入在用
│   │   │       │   └── resource/          #  資源
│   │   │       ├── model/                 #  單一數據服務
│   │   │       ├── logic/                 #  複雜數據服務
│   │   │       ├── game/                  #  遊戲邏輯服務
│   │   │       ├── register/              #  認證服務
│   │   │       └── table/                 #  地端上報服務
│   │   └── port/
│   │       └── any/                       #  對應上面每一組的 interface，同樣用「功能」命名
│   │           ├── admin/
│   │           │   ├── authentication/
│   │           │   └── resource/
│   │           ├── model/
│   │           └── logic/
│   │
│   ├── output/                    # 輸出端（driven adapter）：實作 + 端口介面
│   │   ├── application/
│   │   │   ├── mysql/
│   │   │   │   └── model/         #   
│   │   │   │                      #   
│   │   │   ├── resource/
│   │   │   │   ├── model/         #   AdminUserRepository，透過 gRPC ResourceClient 呼叫 resource 服務
│   │   │   │   │                  #   （http 登入用這份，不直接連 DB）
│   │   │   │   └── logic/         #   AppUserRepository
│   │   │   └── producer/
│   │   │       └── model/         #   AMQP UserProducer
│   │   └── port/
│   │       └── any/               #  跟 usecase/port/any 一樣的命名邏輯：這裡的「介面」不分誰在用，
│   │           ├── model/         #   UserRepository / AdminUserRepository / AppUserRepository
│   │           └── logic/         #   AppUserRepository
│   │
│   ├── register/                  # 組裝層：把 container 生好的 handler 註冊到對應的 server/router
│   │                                #   （grpc.RegisterXxxServer / gin.Group / cron.AddFunc ...），
│   │                                #   cmd/ 只管呼叫 XxxInit 拿到 server 物件再 Serve，不碰組裝細節
│   │
│   └── container/                 # wire 組裝根：wire.go 手寫、wire_gen.go 自動產生，別手改後者
│       （每個服務各自一個 Container + InitXxxContainer：FacadeContainer / ResourceContainer /
│        HttpContainer / ConsumerContainer / CronContainer / WebsocketContainer /
│        ClientContainer / CommandContainer）
│
├── pkg/                            # 跟 domain 無關、可重用的通用元件（logger / router / cache / response / aop 等泛用工具）
│
├── config/                         # viper 讀取的 yaml 設定檔，
│
├── proto/                          # protobuf 原始定義（facade/ 對外、resource/ 資料服務、client/ 外部訂閱）
└── pb/                             # protoc 產生的程式碼，對應 proto/ 底下的定義
```


## DI 依賴注入樹狀圖（ResourceContainer）

說明：`A --> B` 代表 A 被注入到 B（A 是 B 的建構依賴），ResourceContainer 為最底層、最終組裝出來的容器。

```mermaid
graph TD
    subgraph bootstrap
    end

    subgraph internal/resource
        Helpers --> MysqlRepositories

        subgraph internal/output
            MysqlRepositories
            MemoryRepositories
            ProducerRepositories
        end

        Helpers --> Usecases
        MysqlRepositories --> Usecases
        Usecases --> GrpcHandlers
        Helpers --> GrpcHandlers

        Interceptors
        Helpers --> Interceptors

        GrpcHandlers --> ResourceContainer
        Interceptors --> ResourceContainer
    end

    bootstrap --> internal/resource
```

文字版（由下往上）：
```
┌ bootstrap ──┐
│             ├─────┐
└──────┬──────┘     │
       ▼            │
┌ pkg ────────┐     │
│             │     │
└──────┬──────┘     │
       │            │
       │            │
       │            │
       ▼            ▼
┌ internal/resource ──────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                    ┌─────────┐                                                      │
│    ┌───────────────────────────────────────────────┤ Helpers ├───────────────────────────────────────────────────┐  │
│    │                                               └────┬────┘                                                   │  │
│    │                                                    │                                                        │  │
│    │                                                    │                                                        │  │
│    │  ┌┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┐                                 │                                                        │  │
│    │  ┆               ┆                                 ▼                                                        │  │
│    │  ┆           ┌ internal/output ──────────────────────────────────────────────────────────┐                  │  │
│    │  ┆           │ ┌─────────────────────┐  ┌─────────────────────┐  ┌─────────────────────┐ │                  │  │
│    │  ┆           │ │  Mysql/Reposities   │  │  Memory/Reposities  │  │ Producer/Reposities │ │                  │  │
│    │  ┆           │ └──────────┬──────────┘  └──────────┬──────────┘  └──────────┬──────────┘ │                  │  │
│    │  ┆           └───────────────────────────────────────────────────────────────────────────┘                  │  │
│    │  ┆               ▲                                 │                                                        │  │
│    │  └┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┘                                 ▼                                                        │  │
│    │                                              ┌───────────┐                                                  │  │
│    │                                              │  Usecases │◀─────────────────────────────────────────────────┘  │
│    │                                              └─────┬─────┘                                                     │
│    │                                                    │                                                           │
│    │                                                    │                                                           │
│    │                                                    ▼                                                           │
│    └───────────────────────────────────────────┬────────┴──────────┐                                                │
│                                                │                   │                                                │
│                                                ▼                   ▼                                                │
│                                      ┌─────────────────┐  ┌──────────────────┐                                      │
│                                      │   GrpcHandlers  │  │   Interceptors   │                                      │
│                                      └────────┬────────┘  └─────┬────────────┘                                      │
│                                               │                 │                                                   │
│                                               └────────┬────────┘                                                   │
│                                                        ▼                                                            │
│                                              ┌─────────────────────┐                                                │
│                                              │  ResourceContainer  │                                                │
│                                              └─────────────────────┘                                                │
└─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘




```



## 服務拓樸

- **facade**：對外 gRPC 入口，`register` / `table` 兩個 handler 目前是 stub，還沒接 usecase。
- **resource**：內部資料服務，直接讀寫 mysql，僅供 facade / http 呼叫（`AdminUserUsecase` 專屬於這條路徑，走 `usecase/application/any/model`）。
- **http**：Gin REST API，目前只有 `Admin/Authentication/Authenticator/SignIn` 這個登入端點，走專屬的 `AuthenticatorUsecase`（`usecase/application/any/admin/authentication`），透過 gRPC 呼叫 resource 服務查帳號。
- **cron / consumer / command**：三個週邊輸入來源，各自觸發「幫 AppUser 加餘額」，共用同一份 `AppUserUsecase.IncreaseBalance`（`usecase/application/any/admin/resource`），底層都是直接打 mysql。三者差異只在參數怎麼來：cron 排程寫死、consumer 解 MQ payload、command 吃 CLI flag。
- **websocket / client**：容器都還在，但目前沒有掛任何 handler，指令啟動後不會處理任何請求（是之前拿掉舊 `User` 垂直切面後留下的空殼，等有新需求再補）。

依賴方向永遠是「外層指向內層」：`input adapter → usecase/port → usecase → output/port ← output adapter`，
`usecase` 完全不知道自己被 http 還是 grpc 還是 cron 呼叫，也不知道資料到底存在 mysql 還是走 gRPC 轉發。


## 如何 watch 开发
1. go mod 安装下载 air 套件
```zsh
go install github.com/air-verse/air@latest
```


## 代碼開發流程(以 Http 為例子)
1. input：建立新的 Http handlers，註冊到 container，container 再註冊到路由上
2. usecase：建立新的 usecase/application/any/<功能> + usecase/port/any/<功能>，把實作註冊到 container，修改 Http handlers 讓 usecase（port）注入
   - 命名用「功能」不是用「adapter 名稱」：如果這個 usecase 未來可能被其他 adapter（cron/consumer/command...）共用，直接放 any 底下就好，不用每個 adapter 各生一份
3. output：建立新的 output/application/mysql/model + output/port/any/model，把實作註冊到 container，修改 usecase 讓 repo-port 注入
*. 如果 container 首次增加 mysql 需要注入 bootstrap.NewMysql 