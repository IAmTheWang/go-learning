# 学习问答笔记(Go,面向 TS/React/Vue 背景)

这份文件不是练习教程(练习教程在各目录的 `notes.md` 里),是**问答复盘笔记**——
把每次跟 AI 深挖"为什么"的内容沉淀下来,方便以后快速复习,不用把整个对话翻一遍。
按主题分类,不按时间顺序,新的一轮问答直接往对应主题下面追加即可。

---

## 1. Go 程序的最基本骨架

```go
package main       // 声明这个文件属于哪个包，每个 .go 文件都必须有，没有例外
import "fmt"       // 引入标准库的 fmt 包（格式化输入输出，类似 console.log）

func main() {      // 程序入口
	fmt.Println("hello")
}
```

- Go 没有"裸脚本"概念——JS 能直接 `node foo.js` 跑,是因为 JS 诞生于浏览器场景,
  设计目标是"随便写几行就能跑"；Go 从第一天就是给大型后端系统设计的编译型语言,
  编译器必须知道每个文件属于哪个包才能做依赖分析。
- `package main` 是**工具链认的特殊包名**(不是关键字,但受 Go 规范约束):
  只有包名恰好是 `main`、且内部有 `func main()`,`go build`/`go run` 才会把它当成
  可执行程序的入口。换成别的包名就只是一个库,不能直接跑。
- **同一个目录下的所有 `.go` 文件必须声明同一个包名**——这是我们踩过的真实坑:
  在 `05-goroutines-channels/`(已经是 `package concurrency`)里新建了一个
  `package main` 的文件,直接导致整个目录编译失败:
  `found packages concurrency (concurrency.go) and main (test.go) in ...`
  解决办法:独立的实验代码统一放进项目根目录新建的 `playground/` 里,
  跑法:`go run playground/文件名.go`。

### 引号:单引号 ≠ 双引号(和 JS 完全不同)

| 写法 | 含义 |
|---|---|
| `"..."` | 字符串(支持转义,如 `\n`) |
| `` `...` `` | 原始字符串(不转义,可跨行) |
| `'a'` | **单个字符**(rune,本质是数字),不是字符串;`'hello'` 是编译错误 |

---

## 2. 变量、赋值、零值

### `:=` vs `=`

```go
x := 5   // 声明 + 赋值（第一次），类型自动推断
x = 10   // 纯赋值，前提是 x 已经存在
x := 20  // 报错！同一作用域内不能对已声明变量再用 :=
```
对照 JS:`:=` ≈ `let x = 5`,`=` ≈ 后面单独的 `x = 10`。

### 零值(Go 独有的重要特性)

**变量只要声明了,哪怕没赋值,也一定有一个确定的初始值**,不会像 JS 的
`undefined` 那样处于"未定义"状态:

```
int:            0
float64:        0
bool:           false
string:         ""
*int (指针):     <nil>
[]int (slice):  []
map:            map[]
chan:           <nil>
```

指针、slice、map、chan、func、interface 这几种引用类语义的零值都是 `nil`。
**注意**:`nil` slice 可以安全读(`len`、`range`、甚至 `append` 都没问题),
但 `nil` map **不能写**(赋值直接 panic:`assignment to entry in nil map`),
`nil` channel 收发会永久阻塞。这也是为什么 slice/map/channel 需要 `make`
才能变成"真正能用"的状态——它们的零值虽然不报错,但并不完整可用。

### `_`(空白标识符)

`_` 只能"被写入",不能"被读出"——它是"这个位置必须填东西,但我不要"的占位符,
不是"先声明再删除"。用来丢弃你不关心的多返回值:

```go
v, ok := <-ch   // ok 表示 channel 是否还开着
v, _ := <-ch    // 不关心 ok，从一开始就不给它起名字
```

---

## 3. 指针(`&` / `*`)——Go 里最反直觉的部分

### 核心比喻:储物柜

内存是一整面编号储物柜的墙。`x := 42` 是把 42 放进某个柜子(比如 5020 号)。
`p := &x` **不是**复制一份 42,而是造一张写着"5020"的纸条——`p` 存的是**柜子编号**,
不是值本身。`*p` 就是"拿着纸条走到 5020 号柜,把东西取出来"。

- `&x`:值 → 指针(取地址)
- `*p`:指针 → 值(解引用)

两者是**反方向**的操作,不是同一件事;`*` 还有第二个身份——写在类型位置
(`var p *int`)是"声明指针类型",写在表达式位置(`*p`)是"解引用操作",
同一个符号、两种语境。

### 实测:普通赋值 vs 指针,谁会互相影响

```go
x := 42
y := x   // 拷贝，y 跟 x 从此无关
p := &x  // 指针，p 指向 x 那块内存

y = 100
// x = 42  y = 100   (x 没变)

*p = 100
// x = 100  *p = 100  (x 也变了，因为是同一块内存)
```

### 为啥后端要"玩指针"(前端是隐式的,Go 是显式的)

JS 里对象/数组天生"传引用",语言帮你隐式决定;Go 的哲学是"显式优于隐式"——
要不要共享同一块内存,自己用 `&`/`*` 写出来。最典型的动机:**让函数能修改
调用者的变量**(Go 函数参数默认是拷贝传递):

```go
func doubleValue(n int)  { n = n * 2 }     // 只改了拷贝，调用者的 x 不变
func doublePointer(n *int) { *n = *n * 2 } // 改的是地址指向的真身，调用者的 x 会变
```

---

## 4. slice 与 map

- **slice** ≈ TS/JS 的 `Array`,是"能增长的数组"类型,可以 `append`。
- **map** ≈ TS/JS 的 `Map`(行为上更像 Map,不像 `Record`——是真正的数据结构,
  不是编译期类型标注)。前端很少主动用 `Map`,因为大部分场景图省事用 `{}`;
  Go 没有字面量简写,必须显式声明 `map[string]int` 这种"键类型+值类型"。
- 入门阶段可以先粗暴地"slice 当 array 用,map 当 Map 用",以后再补两个细节差异:
  ① slice 可能共享同一块底层数组(不像 JS array 天然独立);
  ② Go map 遍历顺序是故意打乱的(JS Map 保证插入顺序)。
- `[]int`、`map[string]int` 本身就是**类型名**(复合类型),不是"数组加 int"
  两个东西拼起来,对应 TS 的 `number[]`。

### struct

把几个有名字的字段打包成一个新类型,类似 TS 的 `interface`/`type` 描述的形状,
但 struct 是真实存在的值,不只是类型标注。可以嵌套(struct 当另一个 struct 的字段类型)。

```go
type Point struct{ X, Y int }
type Circle struct {
	Center Point
	Radius float64
}
```

---

## 5. `make` vs `new`——功能上不是"都是新建一个变量"

| | `new(T)` | `make(T, ...)` |
|---|---|---|
| 适用范围 | 任意类型 | 只能是 slice / map / channel |
| 返回什么 | 指针 `*T`,指向一块**零值**内存 | T 本身(不是指针),而且是**真正初始化好、能直接用**的值 |
| 能否自己实现 | **能**(用泛型 `func MyNew[T any]() *T { var zero T; return &zero }`) | **不能**,是编译器特殊处理,直接调用 runtime 内部函数(`makeslice`/`makemap`/`makechan`) |

实测证明 `new` 对 slice/map 没啥用——`new([]int)` 解引用出来还是 `nil`,
没有被真正初始化,必须用 `make`:

```
new(int) 解引用: 0
*p2 == nil: true   —— new 并没有让它变成可用的 slice
```

**为啥普通代码写不出自己的 `make`**:它需要根据"这是 slice 还是 map 还是 channel"
执行完全不同的逻辑,还要摸 runtime 内部的数据结构(map 的哈希桶、channel 的
环形缓冲区+等待队列),这些不是 Go 语言里能直接构造出来的公开类型。

**`make` 和 `let`/`var`/`const` 不是一回事**:`let`/`var` 管的是"给一个名字
绑定值"(声明层面),`make` 管的是"造一个初始化好的值"(构造层面),对应 TS 里
`new Map()`、`new Array()` 这种构造函数调用,不是 `let`。两者可以同时出现:

```go
s := make([]int, 0)
// ^^ let 层面：声明 s 这个名字
//      ^^^^^^^^^^^^^^^ 构造层面：造一个空 slice
```

### 泛型速览

`func MyNew[T any]() *T` 里,`[T any]` 是**类型参数**列表(`any` 是约束,
意思是随便什么类型都行),调用时 `MyNew[Point]()` 显式指定 `T = Point`,
跟 TS 的 `myNew<Point>()` 写法基本对应。

---

## 6. 为啥内置类型不能定义方法

```go
func (i int) Double() int { return i * 2 }
```
真实编译报错:
```
cannot define new methods on non-local type int
```
Go 规定**只能给自己包里定义的类型加方法**,`int` 不属于任何用户包。这是故意设计,
为了避免 JS 里"猴子补丁"(`Array.prototype.foo = ...`)那种全局副作用问题。
想要自定义行为,得先定义一个新类型:

```go
type MyInt int
func (i MyInt) Double() MyInt { return i * 2 } // 这样可以
```

---

## 7. goroutine / channel / select

### channel 到底是什么

Go 运行时提供的一个**线程安全的队列**,专门用于 goroutine 之间传数据,并且
自带同步能力(没准备好的时候会让操作自动阻塞等待)。

```go
ch := make(chan int) // 无缓冲 channel
ch <- 42              // 发送：箭头指向 ch，没人接收就一直卡在这一行
v := <-ch             // 接收：箭头离开 ch，没人发送就一直卡在这一行
```
记忆法:**箭头永远指向数据流动的方向**。

### 阻塞(blocking)的判断标准——不是按时长,是"会不会等"

阻塞 = 这一行执行不下去,程序卡住直到条件满足,不管等 1 纳秒还是 10 分钟,
性质都一样;非阻塞 = 压根不等,立刻有结果就走,没结果也立刻走(比如 `default`)。

实测对比:
```
[无 default] 等了 201.886625ms 才拿到 1     ← 阻塞：老实等
[有 default] 等了 542ns 就直接走 default 了  ← 非阻塞：不等
```

### select 是什么(TS 没有对应语法)

`select` 是 Go 的语言级控制流,专门用来**同时等待多个 channel 操作**,
类似 `Promise.race()`,但它是阻塞的控制流语句,不是返回值的函数。

- 没有 `default`:阻塞,等到某个 case 就绪才继续。
- 有 `default`:变成**非阻塞检查**——查一眼有没有,没有就立刻走 default,
  绝不等待。这一点 JS 没有对应写法(Promise 没法"看一眼有没有 resolve")。
- **select 只选一次**:多个 case 同时就绪时只执行最先到的那一个,不会全部执行,
  跟 `switch` 一样"选中一个分支,其余不执行"。
- select 本身不会"处理完一次自动继续监听",想要持续监听要自己套 `for`:

```go
for {
	select {
	case v, ok := <-ch:
		if !ok { return } // channel 关闭，退出循环
		fmt.Println("收到:", v)
	}
}
```
这里只有 1 个 goroutine 在反复执行 select,`for` 不会创建新的 goroutine
(只有 `go` 关键字才会)。

### goroutine vs JS async/await 的本质区别

JS 的异步本质还是单线程——同一时刻只有一段代码在执行,等 IO 的时候才能腾出手
去干别的事。goroutine 是能真正跑在多个 CPU 核心上、同时执行的,由 Go 运行时
自己调度到少数几个系统线程上,启动成本比 JS 的 Promise 还便宜。

### WaitGroup / Mutex

- `sync.WaitGroup`(Add/Done/Wait)≈ `Promise.all`(但没有返回值收集)。
- `sync.Mutex` 用来保护多个 goroutine 同时读写同一个共享变量,防止竞态条件
  (race condition)。Go 的哲学是"share memory by communicating"(用 channel
  传数据),而不是"communicate by sharing memory"(锁着共享变量互相抢)——
  Mutex 属于后一种,是在确实需要共享状态时的退而求其次方案。
- `go test -race` / `go run -race` 能检测出真实的竞态问题。

---

## 8. 垃圾回收(GC)对比

- GC = Garbage Collection(垃圾回收)。
- JS 引擎(V8/SpiderMonkey/JavaScriptCore)各自独立实现 GC,ECMAScript 规范
  不规定具体算法。V8 是分代 GC(新生代/老生代),**按内存压力/阈值触发,不是
  按固定时间间隔扫描**。
- Java(JVM)有多种可选垃圾收集器,G1 是默认,高度可调。
- Go 是并发三色标记清除,**非分代设计**,配合逃逸分析(编译期决定变量分配在
  栈上还是堆上)。
- 三者的 GC 逻辑**不一样**,只是解决的问题(自动管理内存,不用手动 free)是同一个。

---

## 9. 其他零散但重要的纠偏

- JS 不是"一切都是引用类型"——基本类型(number/string/boolean/null/undefined/
  symbol/bigint)是值类型,对象/数组/函数才是引用类型。
- `new(T)` 的零值有时不能直接用,TS 里没有完全对应的概念,但 JS 的
  `let arr; arr.push(1)` 会因为 `undefined` 直接报错崩溃——本质上也是
  "声明了但没构造好的东西不能用"的同类问题,只是崩得更早更彻底。
- 心得类比:**TS 像自动挡,Go 像手动挡**——自动挡帮你处理内存共享/复制/回收
  的细节,手动挡把这些操作交到你手里(指针、值/引用的显式选择),上手门槛更高,
  但控制更精细、行为更可预测。

---

## 10. 学习方法论小结

视频教程学编程效率偏低,不如直接用 AI 深挖一个真实项目——**按需追问、
追到懂为止**,比被动看视频强得多。但代价是覆盖面靠自己是否想到要问;
结构化的教程/文档的价值在于提前告诉你"这个领域该知道哪些东西",而不是
等踩到坑才想起来问。更实际的组合:AI 负责讲透,repo 里各目录的 `notes.md`
负责兜底"别漏项",两者配合,而不是二选一。

---

## 踩坑记录(真实发生过的错误,留作参考)

1. **同目录混用不同 package 名**:在 `05-goroutines-channels/` 里新建
   `package main` 文件,导致整个目录 build 失败。→ 独立实验代码放 `playground/`。
2. **`select_demo.go` 报 `undefined: ch`**:排查发现是 `ch := make(chan int)`
   那一行被误注释掉了,变量根本没声明,跟 channel 语法本身无关。
3. **`TestNewCounter` 早期版本会 panic**:`NewCounter()` 返回 `nil` 时直接调用
   `counter()`,触发"调用 nil 函数值"的运行时 panic(而不是干净的测试失败)。
   修复:调用前先 `if counter == nil { t.Fatal(...) }`。
