# 07 - 自定义 error 与 Stringer 接口

## 核心概念

- **自定义 error 类型**:JS 里自定义错误要 `class Foo extends Error`,靠继承。
  Go 没有继承,规则更直白:**任何类型只要有 `Error() string` 这个方法,
  它就是 `error`**——跟第 03 批"接口靠方法签名满足,不用 `implements`"是
  同一套逻辑:

  ```go
  type NotFoundError struct{ ID int }
  func (e *NotFoundError) Error() string {
      return fmt.Sprintf("item %d not found", e.ID)
  }
  ```

  之所以用**指针接收者**(`*NotFoundError` 而不是 `NotFoundError`)是社区
  惯例:error 值通常不需要拷贝语义,指针接收者也让后面 `errors.As` 的用法
  更统一(`errors.As` 期望目标是指针的指针)。

- **`fmt.Errorf` + `%w`:包装 error,而不是重新造一个**:

  ```go
  if err != nil {
      return fmt.Errorf("load config: %w", err)
  }
  ```

  `%w` 是专门的"包装"占位符(区别于 `%v`/`%s`)——用它包出来的新 error,
  内部仍然**链接着**原始 error,不是把原始信息拍扁成一段文本就完事。
  对照 JS:`throw new Error("load config", { cause: err })`(`Error.cause`
  是 ES2022 才有的新特性,做的是同一件事:在抛出新错误时保留原始错误的引用)。

- **`errors.As`:顺着包装链条,找出你要的具体类型**:

  ```go
  var nf *NotFoundError
  if errors.As(err, &nf) {
      fmt.Println("missing ID:", nf.ID)
  }
  ```

  即使 `err` 已经被 `LoadItemConfig` 用 `%w` 包了一层,`errors.As` 依然能
  "看透"这层包装,找到链条深处真正的 `*NotFoundError`,并把它赋值给 `nf`。
  普通的类型断言(`err.(*NotFoundError)`)做不到这一点——一旦被包装过,
  断言会直接失败,因为 `err` 表面上的类型已经变成了 `fmt.Errorf` 内部生成
  的 wrapError,不再是 `*NotFoundError` 本身。

  `&nf` 传的是"指向 `*NotFoundError` 变量的指针",这样 `errors.As` 才能
  在找到匹配类型时把结果**写回**给你,跟 `strconv.Atoi` 系列"多返回值"
  不是一回事,而是"传指针进去被赋值"的经典 Go 模式(见 `LEARNING-LOG.md`
  第 3 节的储物柜比喻)。

- **`errors.Is` vs `errors.As`(这批没直接用到 `errors.Is`,但顺带记录
  区别,以后会遇到)**:
  - `errors.Is(err, target)`——判断"这条链条上有没有**这个具体的 error 值**"
    (比如标准库常见的哨兵 error:`sql.ErrNoRows`)。
  - `errors.As(err, &target)`——判断"这条链条上有没有**这个类型的 error**",
    并把它取出来用(本批用的就是这种,因为我们要拿到 `nf.ID` 这个字段)。

- **`fmt.Stringer` 接口:让 `fmt` 自动帮你转成人类可读的样子**:

  ```go
  type Temperature float64
  func (t Temperature) String() string {
      return fmt.Sprintf("%.1f°C", float64(t))
  }
  ```

  `fmt.Stringer` 是标准库定义的接口,只有一个方法:`String() string`。
  一旦某个类型实现了它,`fmt.Println(t)`、`fmt.Sprintf("%v", t)` 都会
  **自动调用** `t.String()`,而不是打印裸浮点数。这跟 JS 的 `toString()`
  几乎是同一个心智模型:

  ```js
  class Temperature {
    toString() { return `${this.value.toFixed(1)}°C`; }
  }
  `${temp}`  // 模板字符串会自动调用 toString()
  ```

  区别在于:JS 的 `toString()` 是**所有对象天生自带**的方法(继承自
  `Object.prototype`),你只是覆盖它;Go 的 `String()` 是你**主动实现**
  一个标准库认识的接口,类型默认根本没有这个方法。

## 本节任务

在 `errors_interfaces.go` 里实现:`(*NotFoundError).Error`、`FindItem`、
`LoadItemConfig`、`IsNotFound`、`(Temperature).String`。

```bash
go test ./07-errors-interfaces/...
```
