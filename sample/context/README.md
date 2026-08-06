

## 1. context 用途- 取消
```golang

package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("worker stopped:", ctx.Err())
			return
		default:
			fmt.Println("working...")
			time.Sleep(time.Second)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go worker(ctx)

	time.Sleep(3 * time.Second)

	cancel() // 通知停止

	time.Sleep(time.Second)
}
```


## 2. context 用途- Timeout
```golang
package main

import (
	"context"
	"fmt"
	"time"
)

func queryDB(ctx context.Context) error {

	select {

	case <-time.After(5 * time.Second):
		fmt.Println("query success")
		return nil

	case <-ctx.Done():
		return ctx.Err()

	}
}

func main() {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		2*time.Second,
	)

	defer cancel()

	err := queryDB(ctx)

	fmt.Println(err)
}
```



## 3. context 用途- Context 傳遞
```golang
package main

import (
	"context"
	"fmt"
)

func Repository(ctx context.Context) {
	fmt.Println("Repository")
}

func UseCase(ctx context.Context) {
	fmt.Println("UseCase")
	Repository(ctx)
}

func Controller(ctx context.Context) {
	fmt.Println("Controller")
	UseCase(ctx)
}

func main() {

	ctx := context.Background()

	Controller(ctx)

}
```

## 4. context 用途- WithValue

```golang
package main

import (
	"context"
	"fmt"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

func Repository(ctx context.Context) {

	rid := ctx.Value(RequestIDKey)

	fmt.Println(rid)

}

func main() {

	ctx := context.WithValue(
		context.Background(),
		RequestIDKey,
		"abc123",
	)

	Repository(ctx)

}
```

## 5. context 用途- Server
```golang
func Controller(ctx context.Context) error {

	return usecase.Login(ctx)

}

func Login(ctx context.Context) error {

	return repository.Login(ctx)

}

func Login(ctx context.Context) error {

	return db.QueryContext(ctx, sql)

}

```

## 5. context 用途- 停止 goruntine
```golang

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker(ctx context.Context, wg *sync.WaitGroup, name string) {
	defer wg.Done()

	for {
		select {

		case <-ctx.Done():
			fmt.Printf("[%s] stop: %v\n", name, ctx.Err())
			return

		default:
			fmt.Printf("[%s] working...\n", name)
			time.Sleep(time.Second)
		}
	}
}

func main() {

	// 建立一個可取消的 Context
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	// 啟動三個 Goroutine
	wg.Add(3)

	go worker(ctx, &wg, "MySQL")
	go worker(ctx, &wg, "Redis")
	go worker(ctx, &wg, "RPC")

	// 模擬主程式跑 5 秒
	time.Sleep(5 * time.Second)

	fmt.Println("====== Cancel Context ======")

	// 通知全部停止
	cancel()

	// 等待所有 Goroutine 結束
	wg.Wait()

	fmt.Println("All worker stopped.")
}
```