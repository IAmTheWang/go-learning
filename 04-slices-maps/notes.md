# 04 - Slices & Maps

## 和 JS/TS 的主要差异

- **slice 是动态数组**，日常写代码时基本只用 slice（`[]int`），不直接用固定
  长度的 array。行为上大致对应 JS 的 `Array`，但要注意：多个 slice 可能共
  享同一块底层数组，`append` 在容量不够时会分配新数组、不够时可能就地修改
  ——所以「返回一个新 slice、不修改传入的参数」需要自己小心处理（比如先
  `make` 一个新的再拷贝，而不是在原 slice 上直接改)。
- **map 大致对应 JS 的 `Object`/`Map`**：`map[string]int{}` 类似
  `Record<string, number>`。**坑点**：`var m map[string]int` 声明出来的是
  `nil` map，可以读（返回零值）但**写入会直接 panic**，必须先
  `m := make(map[string]int)` 或用字面量 `map[string]int{}` 初始化。
- **没有内建 Set 类型**：常见写法是用 `map[T]bool` 或 `map[T]struct{}` 当
  集合，靠 key 的唯一性去重（`struct{}` 不占内存，是更地道的写法，但这里
  先用 `bool` 更直观）。
- **`range` 遍历**：`for i, v := range slice` 对应 JS 的
  `for (const [i, v] of arr.entries())`；`for k, v := range m` 对应
  `for (const [k, v] of Object.entries(obj))`，但**map 的遍历顺序是随机
  的**，不像 JS 对象会保留插入顺序。

- **`nil` slice vs `make` 出来的空 slice，序列化成 JSON 时不一样**：这在写
  HTTP API 时是个真实的坑（对应 06-http-basics）。
  - 方式 A：只声明、不初始化 —— `var list []string`，此时 `list == nil`
    （底层没有分配数组内存）。`json.Marshal(list)` 得到的是 `null`，前端
    拿到后执行 `list.length` / `list.map()` 会直接报错崩掉。
  - 方式 B：用 `make([]T, 0)` 显式初始化 —— `list := make([]string, 0)`。
    `make` 是 Go 内建函数，给 slice/map/channel 分配内存并初始化；这里两
    个参数分别是元素类型 `[]T` 和初始长度 `0`。此时 `list != nil`，是一
    个真实存在、内容为空的 slice。`json.Marshal(list)` 得到的是 `[]`
    （空数组），前端 `list.length` 拿到 `0`，`list.map()` 也不会报错。
  - **结论**：只要这个 slice 会被序列化成 JSON 返回给前端，就要用
    `make([]T, 0)`（或字面量 `[]T{}`）显式初始化，不要只声明不初始化。

## 本节任务

在 `collections.go` 里实现：`Reverse`、`Dedup`、`WordCount`。
跑 `go test ./04-slices-maps/...` 直到全部通过。
