# go-learning

Go 练习仓库，面向有 TS / React / Vue 背景的前端开发者。每个目录是一个独立的
练习包：`xxx.go` 里有若干 `// TODO` 待实现的函数，`xxx_test.go` 是对应的测试。

## 怎么做

1. 打开某个目录下的 `notes.md`，看这一批概念和 JS/TS 的对应关系与差异。
2. 打开 `*.go`（非 `_test.go` 的那个），把 `// TODO` 的函数体实现掉。
3. 跑测试，看是否通过：

   ```bash
   go test ./01-basics/...
   # 或者一次跑全部
   go test ./...
   ```

4. 全绿之后再进入下一个目录。测试文件不用改，都是照着 Go 的惯例
   （table-driven test）写的，可以当作理解 Go 测试风格的参考。

其他常用命令：

```bash
gofmt -l .        # 检查格式（Go 对格式要求是强制且自动化的，不需要 Prettier 式的讨论）
go vet ./...      # 静态检查常见错误
```

## 练习列表（按顺序做）

| 目录 | 主题 | 对应你已知的概念 |
|---|---|---|
| [01-basics](01-basics/notes.md) | 变量、类型、常量、控制流 | `let`/`const`、`if`/`switch`、无隐式类型转换 |
| [02-functions](02-functions/notes.md) | 多返回值、命名返回值、可变参数、闭包 | 解构返回值、`...args`、closure |
| [03-structs-interfaces](03-structs-interfaces/notes.md) | struct、method、interface | TS 的 `interface`/`type`、结构化类型 |
| [04-slices-maps](04-slices-maps/notes.md) | slice、map | JS 的 `Array`、`Object`/`Map` |
| [05-goroutines-channels](05-goroutines-channels/notes.md) | goroutine、channel、select、WaitGroup、Mutex | `async/await`（单线程）vs 真并行、竞态条件 |

`05` 涉及并发，额外用竞态检测器再确认一遍：

```bash
go test -race ./05-goroutines-channels/...
```

做完这五个之后可以叫我给下一批（error wrapping/自定义 error、简单的 HTTP server）。
