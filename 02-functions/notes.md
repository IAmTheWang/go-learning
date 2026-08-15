# 02 - Functions

## 和 JS/TS 的主要差异

- **多返回值是一等公民**：`func f() (int, error)` 直接返回两个值，不需要
  像 JS 那样包一个数组或对象再解构。调用时 `a, b := f()`。
- **`(value, error)` 是 Go 处理「预期内失败」的标准方式**，取代了 TS 的
  `try/catch`。惯用写法永远是：

  ```go
  result, err := doSomething()
  if err != nil {
      return err // 或者处理它
  }
  ```

  没有异常会「自己往上冒泡」——每一层都必须显式检查并传递 `err`。
  `panic`/`recover` 是有的，但只用于真正的程序错误（数组越界等），不用来做
  常规的错误处理。
- **命名返回值（named return values）**：`func f() (min, max int)` 中
  `min`、`max` 在函数体内当普通变量用，函数末尾写 `return`（不带值）会自动
  返回它们当前的值。
- **可变参数 `...T`**：等价于 JS 的 `...args`，但类型固定，如
  `func f(nums ...int)`；调用时对已有 slice 展开是 `f(slice...)`
  （对应 JS 的 `f(...arr)`）。
- **闭包（closure）语义和 JS 一致**：函数可以捕获并修改外层变量，`NewCounter`
  这个练习和 JS 里 `function makeCounter() { let n = 0; return () => ++n }`
  是同一个模式。

## 本节任务

在 `functions.go` 里实现：`Divide`、`MinMax`、`NewCounter`。
跑 `go test ./02-functions/...` 直到全部通过。
