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

- **`fmt.Errorf` 的本质:拼接变量字符串,变成一个 `error` 对象**。JS 里
  可以直接用模板字符串扔进 `throw new Error(...)`;Go 的 `error` 是接口,
  字符串不能直接当 `error` 用,所以要靠 `fmt.Errorf` 格式化出一个:

  ```go
  userId := 42
  return fmt.Errorf("user %d not found", userId)
  ```

  **`errors.New` vs `fmt.Errorf`**,区别只在"要不要拼变量":

  | 方法 | TS 对应写法 | 用途 |
  | --- | --- | --- |
  | `errors.New("permission denied")` | `new Error("固定文本")` | 纯静态文本,不带变量 |
  | `fmt.Errorf("user %d not found", id)` | `` new Error(`带 ${var} 的文本`) `` | 动态文本,需要拼接变量 |

  常用占位符(这批最常见的四个):`%s`(字符串)、`%d`(整数)、
  `%v`(通用,几乎任何类型都能填)、`%w`(专用于包装 error,见上一条)。

  **占位符怎么知道该匹配哪个变量**:纯粹按**从左到右的位置顺序**一一
  对应,跟类型无关——`%d` 不是因为"认出"`id` 是数字才匹配它,只是因为
  它是第 1 个占位符、`id` 是第 1 个参数:

  ```go
  fmt.Errorf("load item %d failed: %w", id, err)
  //                     ^^          ^^
  //                   第1个占位符   第2个占位符
  //                     ↓            ↓
  //                    id           err   (按位置对应,不按类型)
  ```

  如果类型和占位符对不上(比如把 `%w` 错配给 `id`、把不存在的 `%r` 配
  给 `err`),Go **不会 panic**,而是在结果字符串里打印一段错误提示,
  同时该占位符对应的"效果"完全作废:

  ```go
  fmt.Errorf("load item %w failed: %r", id, err)
  // 结果类似:"load item %!w(int=123) failed: %!r(string=item 42 not found)"
  // - %w 配上了 int 类型的 id → 包装失败,输出 %!w(int=123)
  // - %r 根本不是 fmt 认识的动词 → 输出 %!r(...)
  // 后果:err 没有真正被 %w 包装,后续 errors.Is/errors.As/errors.Unwrap
  // 都找不到原始的 err 了。
  ```

  如果占位符顺序和参数顺序对不上,或同一个参数要重复用,可以用**显式
  索引** `%[n]`,明确指定"用第 n 个参数":

  ```go
  fmt.Errorf("error %[2]w occurred on item %[1]d", id, err)
  // %[1]d → 第 1 个参数 id;%[2]w → 第 2 个参数 err
  ```

  **注意这里的 `n` 是从 1 开始数的**,跟 Go 里数组/切片/字符串从 0 开始
  索引(`arr[0]`)是两套不同的编号系统——`%[n]` 数的是"这是格式化字符串
  后面第几个参数"(第 1 个、第 2 个……),不是内存里的下标,这个习惯来自
  C 的 `printf`(`%1$s`)传统,和正则替换里 `$1`、`$2` 是同一个心智模型。

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

- **`errors.As` 底层原理:本质是一个手写的循环,不是语言内置语法糖**。
  对照自己用 TS 写的等价逻辑:

  ```ts
  function isNotFound(err: Error): boolean {
    let current: any = err;
    while (current) {
      if (current instanceof NotFoundError) return true;
      current = current.cause; // 剥下一层
    }
    return false;
  }
  ```

  `errors.As` 做的事完全一样,只是 Go 标准库里没有 `instanceof`,得靠
  `reflect` 包在运行时比较类型。简化版逻辑大致是:

  ```go
  func As(err error, target any) bool {
      for err != nil {
          if reflect.TypeOf(err) == reflect.TypeOf(target).Elem() {
              reflect.ValueOf(target).Elem().Set(reflect.ValueOf(err))
              return true
          }
          err = errors.Unwrap(err) // 相当于 current.cause
      }
      return false
  }
  ```

  （真实标准库实现更复杂,还处理了类型是否实现自定义 `As(any) bool` 方法
  等情况,但核心骨架就是上面这样。）

  对应关系:

  | TS 版本 | Go 版本 |
  | --- | --- |
  | `current.cause` | `errors.Unwrap(err)`(要求该 error 实现 `Unwrap() error`;`fmt.Errorf("...%w", err)` 生成的 wrapError 自动实现了它) |
  | `current instanceof NotFoundError` | `reflect.TypeOf(err) == reflect.TypeOf(target).Elem()` |
  | `while (current)` | `for err != nil` |
  | 只返回 `bool` | 匹配上时还会把 `err` **写回** `target` 指向的变量,再返回 `true` |

  逐行拆解这行反射代码——`reflect.TypeOf(err) == reflect.TypeOf(target).Elem()`：

  - `reflect.TypeOf(x)`:拿到 `x` 在**运行时**的具体类型,类比 JS 的
    `x.constructor`(比 `typeof x` 精确得多,`typeof` 在 JS 里对对象只会
    统一返回 `"object"`)。
  - 调用处传入的是 `errors.As(err, &target)`,其中
    `var target *NotFoundError`——所以 `&target` 的类型是
    `**NotFoundError`(指向"指向 `NotFoundError` 的指针"的指针)。
  - `.Elem()`:剥掉**一层**指针/引用,拿到里面那层的类型。类比 TS 类型
    层面的 `T extends *infer U ? U : never`:`**NotFoundError` 经过一次
    `.Elem()` 变成 `*NotFoundError`。
  - 整行翻译成人话:"**当前这层 `err` 的具体类型,是不是恰好等于
    `*NotFoundError`?**",约等于 JS 里的 `current.constructor === NotFoundError`。

  为什么要绕这一圈用反射:因为 `As` 函数签名里 `target` 参数类型是
  `any`,函数体内**拿不到你想匹配的类型本身**——你没有直接传
  `*NotFoundError` 这个类型进去,而是传了"指向该类型变量的指针" `&target`。
  函数只能靠反射"打开"这个指针,看它指向的东西是什么类型,才知道调用者
  到底想找哪种 error。这也是为什么必须传指针:**既用来告诉函数目标类型
  是什么,又用来在匹配成功后把结果写回调用者的变量**——一举两得,和
  上面「传指针进去被赋值」是同一个模式的另一种应用。

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

  `fmt.Stringer` 的标准库源码非常简单:

  ```go
  type Stringer interface {
      String() string
  }
  ```

  它的契约只有一条:**任何类型,只要拥有一个名为 `String()`、签名是
  `func() string` 的方法,就自动满足这个接口**——不需要写
  `implements Stringer`,不需要显式声明,`fmt` 包只在运行时检查"这个
  值有没有这个方法",这跟第 03 批"接口靠方法签名满足"的结构化类型
  (structural typing)是同一套规则,`Stringer` 只是标准库里最常见的
  一个应用案例。

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
