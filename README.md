## 六角架構圖

文字版（由下往上）：

```
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

```

## 六角框架使用重點

1. input 不同輸入 adaptor 相同的路徑應該是相同的設備方法, 譬如 
   
    cron/admin/authentication/authenticator_handler.go,  
    command/admin/authentication/authenticator_handler.go,  
    http/admin/authentication/authenticator_handler.go,  
    facade/admin/authentication/authenticator_handler.go,  
  
    以不同協議但是相同業務邏輯實現 登入取得 token 這個業務。
2. output 同理
3. 六角框架的也間接看出，後端的本質其實是在做消息傳遞，不論輸入輸出怎麼變

## 如何利用這套框架 從 1到100 建立一個獨立 grpc服務端-服務

1. proto基本文檔撰寫:
  建立獨立目錄的 proto/source/announcement/lottery.proto, 並生成
2. config設定當基本參數:
  設定 config/services.yaml  與 bootstrap/config.go
3. cmd 撰寫啟動入口程序：
  建立 cmd/source.go 與 container tree 與 register
4. input & handler 協議輸入代碼建立：
  建立 input/source/announcement/lottery.go 的服務類，並且綁定 pb， 並注入 container (記得注入 abstract 類)，並且註冊 register/
5. usecase 業務邏輯建立：
  5-1 建立 usecase 包含 介面跟實作，（實作先簡單return 寫死 domain 數據），  
  5-2. 並且注入 handler 與 container  (記得注入 abstract 類)，並且修改 input 使用 usecase 類
6. output 資料輸出建立：
  6-1 .建立 output 包含 介面跟實作，（實作先簡單return 寫死 domain 數據），  
  6-2. 並且注入 usecase 與 container (記得注入 abstract 類), 並且修改 usecase 使用 repository 類

## 如何 建立一個獨立 grpc客戶端-連線 conn。

1. proto: 建立獨立目錄的 proto/source/announcement/lottery.proto, 並生成
2. config: 設定 config/clients.yaml  與 bootstrap/config.go
3. bootstrap： bootstrap/source.go 基本客戶端連線（conn） + internal/client/source_client.go
4. cmd ：建立 cmd/daemon.go 與 container tree 與 register/daemon.go

## 我們的服務分成兩層

我們把服務定義成兩個層級：

- **第一層是協議層服務**：  
本質上只是一個服務載體，對應不同的協議實作方式（http / grpc / command / cron ...），本身不包含任何業務邏輯。  
- **第二層才是邏輯服務**：真正的業務邏輯放在這一層，例如 `admin` 服務屬於邏輯層，同一份邏輯可以同時掛載到 http、grpc、command 等不同的協議載體上，不綁定特定協議。

1. 協議服務（實例服務）：
  - Source：開獎資料來源服務載體
    - Daemon：常駐任務服務載體
    - Facade：對外主要 gRPC 服務載體
    - Resource：資料 gRPC 服務載體
    - Http：主要 HTTP 服務載體
    - Command：CLI 服務載體
    - Cron：排程服務載體
    - Websocket：WebSocket 服務載體
2. 邏輯服務（虛擬服務）：
  - admin：後台介面邏輯服務，負責所有後台相關業務邏輯
  - app：前台介面邏輯服務，負責所有前台相關業務邏輯
  - third：第三方介接邏輯服務，負責所有第三方串接相關業務邏輯
  - game：前台遊戲介面邏輯服務，負責所有前台遊戲相關業務邏輯
  - table：前台資料介面邏輯服務，負責所有前台地端上報相關業務邏輯
  - register：前台驗證邏輯服務，負責所有前台身份驗證相關業務邏輯
  - logic：次級（衍生）資料邏輯服務，處理跨多個資源、需要額外組合運算的資料邏輯
  - model：次級資料的增刪改查（CRUD）邏輯服務
  - announcement：開獎邏輯服務
  - watcher：採集開獎資料的邏輯服務
3. 實際運作服務堆疊

```
協議服務（載體）        實際掛載的邏輯服務
+-------------+     +--------------------------------+
|   Source    | --> | announcement                   |
+-------------+     +--------------------------------+
|   Daemon    | --> | watcher                        |
+-------------+     +--------------------------------+
|   Facade    | --> | game, table, register, admin   |
+-------------+     +--------------------------------+
|   Resource  | --> | logic, model                   |
+-------------+     +--------------------------------+
|   Http      | --> | admin, app, third              |
+-------------+     +--------------------------------+
|   Command   | --> | admin                          |
+-------------+     +--------------------------------+
|   Cron      | --> | admin                          |
+-------------+     +--------------------------------+
|   Websocket | --> | (尚未掛載任何邏輯服務)            |
+-------------+     +--------------------------------+
```

## 傳統 ThinkPHP MVC 遇到的核心問題

1. ThinkPHP 裡面的 Controller 雖然是業務主邏輯，但是裡面包含了 http 協議的代碼，這個會造成多 不同輸入實作有困難

## 這套框架的各個職責

1. input : 只是寫協議的對接 （如 grpc http command），決定這個服務要用在哪一個服務載體
2. usecase : 業務邏輯，基本上就是 Tp 的 C 去掉了協議的部分。
3. output/**/model：負責單一數據操作的 數據模型。基本上就是 Tp 的 M-mdoel
  output/**/logic：負責複雜數據操作的 數據模型。基本上就是 Tp 的 M-logic
4. 基本上，就是這四個元件交互

## 目錄結構

```
.
├── main.go                        #   全程式進入點，只呼叫 cmd.Execute()
├── cmd/                           #   cobra 指令，每個檔案對應一個可獨立啟動的服務／進程
│   ├── root.go                    #   root command，Execute() 供 main.go 呼叫
│   ├── facade.go                  #   啟動 facade gRPC 服務（對外入口）
│   ├── resource.go                #   啟動 resource gRPC 服務（資料服務，僅供 facade / http 呼叫）
│   ├── http.go                    #   啟動 HTTP（Gin）服務
│   ├── consumer.go                #   啟動 AMQP consumer
│   ├── client.go / socket.go      #   啟動 gRPC client-side stream 訂閱
│   ├── source.go                  #   啟動開獎資料來源服務（source client 訂閱）
│   ├── daemon.go                  #   啟動常駐任務（watcher）
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
│   │       ├── facade/            #   對外 gRPC 入口：register/、table/、admin/authentication/
│   │       ├── resource/          #   resource 內部 gRPC 服務（僅供 facade / http 呼叫）：model/
│   │       ├── http/              #   對外 Http 入口：admin/authentication/、admin/resource/
│   │       ├── client/            #   gRPC client（訂閱外部 stream）：admin/resource/
│   │       ├── consumer/          #   AMQP consumer：admin/resource/
│   │       ├── source/            #   開獎資料來源服務載體：announcement/
│   │       ├── daemon/            #   常駐任務服務載體：watcher/source/
│   │       ├── cron/              #   排程任務：admin/authentication/、admin/resource/
│   │       ├── websocket/         #   websocket 入口：admin/resource/
│   │       └── command/           #   CLI 指令：admin/authentication/、admin/resource/
│   │       （每個 adapter 底下都有自己獨立的 abstract_handler.go，彼此不共用；
│   │        adapter 內部依 leaf 功能再分 admin/resource、admin/authentication 這種子資料夾——
│   │        這是跨層的對應 key，不是字面語意：input/<adapter>/admin/resource、
│   │        usecase/application/any/admin/resource、usecase/port/any/admin/resource
│   │        三者相對路徑相同，代表同一條 usecase 邏輯，不代表跟 HTTP 後台路由有關；
│   │        目前四個 adapter：http / command / cron / facade 都已經掛上同一條
│   │        admin/authentication 登入邏輯，是這套框架「一個 usecase、多種載體」的主要範例）
│   │
│   ├── middleware/
│   │   └── admin/                    # HTTP 專用 middleware 鏈
│   │
│   ├── interceptor/                  # gRPC 專用攔截器鏈
│   │   ├── facade/
│   │   │   └── game/
│   │   └── resource/
│   │
│   ├── usecase/                           # 商務案例：實作 + 端口介面
│   │   ├── application/                   # 實作
│   │   │   └── any/                       #  「any」表示這份 usecase 不綁定特定 adapter，可以被多個 driving
│   │   │       │                          # usecase端一種命名結構對應不同input 但是多種輸入
│   │   │       ├── admin/                 #  後台服務
│   │   │       │   ├── authentication/    #  登入邏輯，被 http/command/cron/facade 四個 adapter 共用
│   │   │       │   └── resource/          #  資源
│   │   │       ├── model/                 #  單一數據服務
│   │   │       ├── logic/                 #  複雜數據服務
│   │   │       ├── game/                  #  遊戲邏輯服務
│   │   │       ├── register/              #  前台驗證邏輯服務
│   │   │       ├── table/                 #  地端上報服務
│   │   │       ├── annoucement/           #  開獎邏輯服務（資料夾拼字沿用舊名，跟 port/any/announcement 拼法不同步，屬已知瑕疵）
│   │   │       └── watcher/
│   │   │           └── source/            #  採集開獎資料的邏輯服務
│   │   └── port/                          #  介面
│   │       └── any/                       #  對應上面每一組的 interface，同樣用「功能」命名
│   │           ├── admin/
│   │           │   ├── authentication/
│   │           │   └── resource/
│   │           ├── model/
│   │           ├── logic/
│   │           ├── game/
│   │           ├── register/
│   │           ├── table/
│   │           ├── announcement/
│   │           └── watcher/
│   │               └── source/
│   │
│   ├── output/                    # 輸出端（driven adapter）：實作 + 端口介面
│   │   ├── application/           # 實作
│   │   │   ├── mysql/
│   │   │   │   ├── model/         #   AdminUserRepository / AppUserRepository，直連 DB
│   │   │   │   └── logic/         #   複雜查詢邏輯，直連 DB
│   │   │   ├── resource/
│   │   │   │   ├── model/         #   AdminUserRepository，透過 gRPC ResourceClient 呼叫 resource 服務
│   │   │   │   │                  #   （http / facade 登入用這份，不直接連 DB）
│   │   │   │   └── logic/         #   AppUserRepository
│   │   │   ├── cache/
│   │   │   │   ├── model/         #   Redis 直接讀寫，不嵌套/不包裝底層 repository
│   │   │   │   └── logic/
│   │   │   ├── memory/
│   │   │   │   ├── model/         #   進程內記憶體快取
│   │   │   │   └── logic/
│   │   │   └── producer/
│   │   │       └── model/         #   AMQP UserProducer
│   │   └── port/                  #  介面
│   │       └── any/               #  跟 usecase/port/any 一樣的命名邏輯：這裡的「介面」不分誰在用，
│   │           ├── model/         #   UserRepository / AdminUserRepository / AppUserRepository
│   │           └── logic/         #   AppUserRepository
│   │
│   └── register/                  # 組裝層：把 container 生好的 handler 註冊到對應的 server/router
│                                    #   （grpc.RegisterXxxServer / gin.Group / cron.AddFunc ...），
│                                    #   cmd/ 只管呼叫 XxxInit 拿到 server 物件再 Serve，不碰組裝細節
│
├── container/                      # wire 組裝根（跟 internal/ 平級，不在裡面）：
│                                    #   wire.go 手寫、wire_gen.go 自動產生，別手改後者
│       （每個服務各自一個 Container + InitXxxContainer：FacadeContainer / ResourceContainer /
│        HttpContainer / ConsumerContainer / CronContainer / WebsocketContainer /
│        ClientContainer / CommandContainer / SourceContainer / DaemonContainer）
│
├── pkg/                            # 跟 domain 無關、可重用的通用元件（logger / router / cache / response / aop 等泛用工具）
│
├── config/                         # viper 讀取的 yaml 設定檔，
│
├── proto/                          # protobuf 原始定義（facade/ 對外、resource/ 資料服務、client/ 外部訂閱、source/ 資料來源）
└── pb/                             # protoc 產生的程式碼，對應 proto/ 底下的定義
```

## DI 依賴注入樹狀圖（ResourceContainer）

文字版（由下往上）：

```
┌ bootstrap ──┐
│             ├─────┐
└──────┬──────┘     │
       ▼            │
┌ pkg ────────┐     │
│             │     │
└──────┬──────┘     │      
       ▼            ▼
┌ internal/resource ───────────────────────────────────────────────────────────────────────────────────────────┐
│                                               ┌─────────┐                                                    │
│   ┌───────────────────────────────────────────┤ Helpers ├─────────────────────────────────────────────────┐  │
│   │                                           └────┬────┘                                                 │  │
│   │  ┌┄┄┄┄┄┄┄┄┄┄┄┐                                 │                                                      │  │
│   │  ┆           ┆                                 ▼                                                      │  │
│   │  ┆       ┌ internal/output ──────────────────────────────────────────────────────────┐                │  │
│   │  ┆       │ ┌─────────────────────┐  ┌─────────────────────┐  ┌─────────────────────┐ │                │  │
│   │  ┆       │ │  Mysql/Reposities   │  │  Memory/Reposities  │  │ Producer/Reposities │ │                │  │
│   │  ┆       │ └──────────┬──────────┘  └──────────┬──────────┘  └──────────┬──────────┘ │                │  │
│   │  ┆       └───────────────────────────────────────────────────────────────────────────┘                │  │
│   │  ┆           ▲                                 │                                                      │  │
│   │  └┄┄┄┄┄┄┄┄┄┄┄┘                                 ▼                                                      │  │
│   │                                          ┌───────────┐                                                │  │
│   │                                          │ Usecases  │◀───────────────────────────────────────────────┘  │
│   │                                          └─────┬─────┘                                                   │
│   │                                                ▼                                                         │
│   └───────────────────────────────────────┬────────┴──────────┐                                              │
│                                           ▼                   ▼                                              │
│                                 ┌─────────────────┐  ┌──────────────────┐                                    │
│                                 │   GrpcHandlers  │  │   Interceptors   │                                    │
│                                 └────────┬────────┘  └─────┬────────────┘                                    │
│                                          └────────┬────────┘                                                 │
│                                                   ▼                                                          │
│                                        ┌─────────────────────┐                                               │
│                                        │  ResourceContainer  │                                               │
│                                        └─────────────────────┘                                               │
└──────────────────────────────────────────────────────────────────────────────────────────────────────────────┘




```

## 服務拓樸

- **facade**：對外 gRPC 入口，`register` / `table` 兩個 handler 目前是 stub，還沒接 usecase。
- **resource**：內部資料服務，直接讀寫 mysql，僅供 facade / http 呼叫（`AdminUserUsecase` 專屬於這條路徑，走 `usecase/application/any/model`）。
- **http**：Gin REST API，目前只有 `Admin/Authentication/Authenticator/SignIn` 這個登入端點，走專屬的 `AuthenticatorUsecase`（`usecase/application/any/admin/authentication`），透過 gRPC 呼叫 resource 服務查帳號。
- **cron / consumer / command**：三個週邊輸入來源，各自觸發「幫 AppUser 加餘額」，共用同一份 `AppUserUsecase.IncreaseBalance`（`usecase/application/any/admin/resource`），底層都是直接打 mysql。三者差異只在參數怎麼來：cron 排程寫死、consumer 解 MQ payload、command 吃 CLI flag。
- **websocket / client**：容器都還在，但目前沒有掛任何 handler

依賴方向永遠是「外層指向內層」：`input adapter → usecase/port → usecase → output/port ← output adapter`，
`usecase` 完全不知道自己被 http 還是 grpc 還是 cron 呼叫，也不知道資料到底存在 mysql 還是走 gRPC 轉發。

## 服務拓撲圖：Facade / Http / Websocket → Resource → Redis / MySQL；Command 直連 MySQL

`facade` 跟 `http` 都沒有自己的資料庫連線，登入查帳號一律經 gRPC 轉發給 `resource`；`resource` 才是真正碰資料庫（mysql）／快取（redis）的那一層。
（目前 `resource` 的 `bootstrap.NewRedis` 已經接進 `ResourceContainer`，但 `AdminUserUsecase` 這條路徑還只走 mysql，redis 連線已備好、尚未有 usecase 使用；
`websocket` 目前容器是空殼，還沒接 `ResourceClient` 也沒掛任何 handler，這裡先畫上去代表「預期中」會走的路徑）

`command` 則是唯一不經過 `resource` 的協議層：`CommandContainer` 直接 wire 了 `outputApplicationMysqlModel`，
繞過 gRPC 自己連 mysql——跟 `resource` 打的是同一個 MySQL instance，只是省了一趟 gRPC。

```
+-------------+  +-------------+  +-------------+     +-------------+
|    Facade   |  |     Http    |  |  Websocket  |     |   Command   |
+-------------+  +-------------+  +-------------+     +-------------+
       |                |                |                   |
       +----------------+----------------+                   |
                        |                                    |
                        | gRPC 呼叫                          | 直連 mysql
                        v                                    | （不走 resource）
             +--------------------+                          |
             |      Resource      |                          |
             +--------------------+                          |
                        |                                    |
                  +-----+-----+                              |
                  |           |                              |
                  v           v                              |
            +-----------------------+
            | +-------+   +-------+ |
            | | Redis |   | MySQL | |
            | +-------+   +-------+ |
            +-----------------------+

```

## 如何 watch 开发

1. go mod 安装下载 air 套件

```zsh
go install github.com/air-verse/air@latest
```

## 微服務地端如果有強大的數據需求，又對網路延遲鳴感，怎麼辦？

1. 在地端 使用 SQL lite 思維，本地 SQL思維
2. 再統一上報數據

