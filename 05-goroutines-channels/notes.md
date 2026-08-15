# 05 - Goroutines & Channels

## 核心概念

- **goroutine**:在函数调用前加一个 `go`,这段代码就会另起一条执行线,和
  当前代码同时跑,不等它结束就往下走:

  ```go
  go doSomething() // 不阻塞，立刻继续执行下一行
  ```

  和 JS 的 `async/await` 不同——JS 永远只有一条执行线,goroutine 是可以
  真正跑在多个 CPU 核心上的（细节看上一轮对话里的厨房比喻）。

- **channel**:goroutine 之间传数据的管道,类型化、线程安全：

  ```go
  ch := make(chan int)   // 无缓冲 channel
  ch <- 42                // 发送（没人接收就一直阻塞在这一行）
  v := <-ch                // 接收（没人发送就一直阻塞在这一行）
  ```

  `make(chan int, 3)` 是带缓冲的 channel，缓冲区满之前发送不阻塞。

- **`select`**：同时等待多个 channel，哪个先有动静就走哪个分支。常见用法
  是配合 `time.After` 做超时保护，避免代码永远卡死：

  ```go
  select {
  case v := <-ch:
      // 正常收到结果
  case <-time.After(time.Second):
      // 一秒内没收到，超时处理
  }
  ```

- **`sync.WaitGroup`**：等待一批 goroutine 全部跑完，等价于 JS 里的
  `Promise.all(...)`，但没有返回值，只是单纯的"等大家都结束"：

  ```go
  var wg sync.WaitGroup
  for i := 0; i < 5; i++ {
      wg.Add(1)          // 登记一个待完成任务
      go func() {
          defer wg.Done() // 完成时打卡
          // ... 干活 ...
      }()
  }
  wg.Wait() // 阻塞直到所有 goroutine 都 Done() 过
  ```

- **`sync.Mutex`**：Go 的哲学是"通过通信共享内存"(用 channel)，但很多时候
  用锁保护一个共享变量更直接。多个 goroutine 同时读写同一个变量而不加保护，
  就是"竞态条件（race condition）"——这正是上一轮说的"前端很少遇到，后端
  必须理解"的并发坑。

  ```go
  type Counter struct {
      mu    sync.Mutex
      count int
  }
  func (c *Counter) Increment() {
      c.mu.Lock()
      defer c.mu.Unlock()
      c.count++
  }
  ```

## 亲眼看一次竞态（可选，不计入测试）

想直观感受"竞态条件"是什么，可以把下面这段粘到一个临时文件里跑跑看：

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	count := 0
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count++ // 危险：多个 goroutine 同时读写同一个变量
		}()
	}
	wg.Wait()
	fmt.Println(count) // 大概率不是 1000
}
```

用 `go run -race main.go` 跑，Go 会明确告诉你"DATA RACE"发生在哪一行。
这就是为什么本节练习里的 `SafeCounter` 必须用 `sync.Mutex` 保护。

## 本节任务

在 `concurrency.go` 里实现：`SquareAsync`、`SumConcurrent`、`SafeCounter`
（含 `Increment`、`Value` 方法）。

```bash
go test ./05-goroutines-channels/...
# 额外用竞态检测器再跑一遍，确认 SafeCounter 真的安全：
go test -race ./05-goroutines-channels/...
```
