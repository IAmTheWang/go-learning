# 06 - net/http Basics

## 核心概念

- **handler 函数**:一个 handler 就是一个签名固定的函数:

  ```go
  func HelloHandler(w http.ResponseWriter, r *http.Request) {
      // 从 r 里读请求信息，往 w 里写响应
  }
  ```

  对照 Express:`app.get('/hello', (req, res) => { ... })`——`r *http.Request`
  相当于 `req`,`w http.ResponseWriter` 相当于 `res`,只是 Go 没有框架帮你
  自动注册路由参数,一切都靠标准库这两个对象手动操作。

- **`http.ResponseWriter` 是接口,不是你 `new` 出来的东西**:真正在生产环境
  跑的时候,Go 的 HTTP server 会自动帮你造一个具体实现传进来;写测试的时候,
  `httptest.NewRecorder()` 造一个"假的"实现,把你写的内容记录下来,方便断言:

  ```go
  rec := httptest.NewRecorder()
  HelloHandler(rec, req)
  rec.Body.String() // 拿到你写进去的内容，直接断言，不需要真的起一个端口
  ```

- **读 query 参数**:`r.URL.Query().Get("name")`——如果 URL 里没有这个参数,
  返回的是空字符串 `""`,不会报错、不会 panic。跟 map 读不存在的 key 是
  同一套"零值,不是异常"的哲学。

- **写响应**:`w.Write([]byte(...))` 写原始字节,或者用 `fmt.Fprintf(w, ...)`
  直接把格式化字符串写进去(`http.ResponseWriter` 本身实现了 `io.Writer`
  接口,所以能直接传给 `fmt.Fprintf`)。

- **设置响应头**:`w.Header().Set("Content-Type", "application/json")`。
  **顺序很重要**——必须在第一次调用 `Write`/`Encode`(写 body)**之前**设置
  header,一旦开始写 body,响应头就已经被"钉死"发出去了,之后再 `Set` 不会
  生效。

- **`encoding/json` 编码响应体**:

  ```go
  type User struct {
      ID   int    `json:"id"`
      Name string `json:"name"`
  }
  json.NewEncoder(w).Encode(User{ID: 1, Name: "Alice"})
  // 直接写进 w，等价于 fmt.Fprintf(w, `{"id":1,"name":"Alice"}`)
  ```

  **struct tag**(`` `json:"id"` ``)是新语法:字段后面反引号里的内容不是
  注释,是给标准库看的"元数据",告诉 `encoding/json` 序列化这个字段时该用
  哪个 key 名。不写 tag 的话,`encoding/json` 默认直接用 Go 的字段名(首字母
  大写),比如 `ID` 会变成 JSON 里的 `"ID"` 而不是 `"id"`——这也是为什么 Go
  的 struct 字段必须首字母大写(exported)才能被 `encoding/json` 看到并
  序列化,小写字段会被直接跳过,一声不吭。

## 本节任务

在 `http_basics.go` 里实现:`HelloHandler`、`UserHandler`。

```bash
go test ./06-http-basics/...
```
