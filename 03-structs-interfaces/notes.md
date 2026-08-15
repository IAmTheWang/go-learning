# 03 - Structs & Interfaces

## 和 JS/TS 的主要差异

- **struct 是纯数据**，类似 TS 的 `type Rectangle = { width: number; height: number }`。
  没有 `class`、没有构造函数、没有继承。
- **方法（method）是「函数 + receiver」**，定义在 struct 外面，而不是像 TS
  的 class 那样写在类体内：

  ```go
  func (r Rectangle) Area() float64 { return r.Width * r.Height }
  ```

  `(r Rectangle)` 叫 receiver，调用时 `rect.Area()`，`r` 就是 `rect` 的一份
  拷贝（value receiver）。如果方法需要修改 struct 本身，receiver 要用指针
  `(r *Rectangle)`，这个先了解概念即可，本节练习都用 value receiver。
- **interface 是结构化类型（structural typing）**——这点和 TS 的 interface
  反而很像：只要一个类型实现了 interface 要求的所有方法，它就自动「满足」
  这个 interface，不需要像 Java/C# 那样写 `implements`。`Rectangle` 和
  `Circle` 都有 `Area() float64`，所以它们都自动是 `Shape`。
- **没有 class 继承，只有组合（embedding）**：Go 用「把一个 struct 塞进
  另一个 struct」来复用字段和方法，本节先不涉及，后面会专门练习。

## 本节任务

在 `shapes.go` 里实现：`Rectangle.Area`、`Circle.Area`、`TotalArea`。
跑 `go test ./03-structs-interfaces/...` 直到全部通过。
