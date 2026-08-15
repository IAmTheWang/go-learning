# 01 - Basics

## 和 JS/TS 的主要差异

- **静态类型，没有 `any`**：变量的类型在编译期确定。`var x int` 或用类型推断
  `x := 5`（`:=` 只能在函数内用，等价于「声明 + 赋值」）。
- **零值，没有 `undefined`**：声明但未赋值的变量会拿到该类型的零值——
  `int` 是 `0`，`string` 是 `""`，`bool` 是 `false`，指针/interface/slice/map
  是 `nil`。没有「未初始化」这种中间状态。
- **没有隐式类型转换**：`1 + 1.0` 编译不过，必须显式 `float64(1) + 1.0`。
  不存在 JS 里 `"1" + 1 === "11"` 这种坑，但你也不能偷懒。
- **只有 `for`**：没有 `while`/`do-while`，`for` 省略掉条件部分就是 `while`：
  `for { ... }` 是无限循环，`for cond { ... }` 相当于 `while (cond)`。
- **`switch` 默认不 fallthrough**：每个 `case` 执行完自动 `break`，和 JS 相反；
  需要贯穿到下一个 case 时要显式写 `fallthrough`。
- **`const` 是真正编译期常量**：不像 JS 的 `const` 只是「不能重新赋值的变量」。

## 本节任务

在 `basics.go` 里实现：`Sum`、`IsEven`、`FizzBuzz`、`Grade`。
跑 `go test ./01-basics/...` 直到全部通过。
