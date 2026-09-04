# DBトランザクション（ROLLBACK）と health-check ハートビート

> 业务知识笔记 — 源自 2026-09-04 实际跑通 welby-fhir-server-aggregator-platform
> 的 floci 本地环境（PR #1505 复审）时，从真实容器日志里挖出来的两个后端基础概念。
> 姊妹篇：[dda-architecture-overview.md](./dda-architecture-overview.md)

## 1. トランザクション（transaction）と ROLLBACK

### 一句话定义

**トランザクション**＝把好几个数据库操作打包成"一件事"，规定这几步**要么全部成功，
要么一步失败就当全部没发生过**——不会出现"改了一半"的中间状态。

可以类比银行柜台办业务：柜员先说"我先帮您锁定这个账户，接下来连续办几步"（= `BEGIN`），
办完了说"都办好了，生效"（= `COMMIT`），中途出问题就说"这次作废，恢复原状"（= `ROLLBACK`）。
不管你中间做的是查询还是修改，只要说了 `BEGIN`，数据库就会维持一个"操作进行中"的状态，
直到你明确 `COMMIT` 或 `ROLLBACK` 为止。

### 为什么一个 SELECT（纯查询）也会用到 ROLLBACK

容易搞混的地方：**ROLLBACK 不是"撤销已经改坏的数据"**——如果这次操作从头到尾都是
`SELECT`（读操作），根本没有数据被改过，谈不上"撤销"什么内容。

真正的原因是"**事务得有始有终，不能开着不管**"。实测遇到的例子（见下面 §3）：
代码逻辑是先开一个事务（`BEGIN`），准备连续查两件事——先数一下"お知らせ"总共
多少条，再去分页取具体内容。第一步的 SQL 就报错失败了，整个"连续要做的两件事"
没法完成，代码必须明确 `ROLLBACK`，把数据库占用的资源（连接、锁）还回去，
不能什么都不说就撒手不管。

**如果不 ROLLBACK 会怎样**：这个"占着但没人管"的事务会一直挂在数据库那边，
时间长了会占用数据库连接池、甚至挡住别的请求的锁，造成比"这一次请求失败"
严重得多的连锁问题（比如整个服务的数据库连接被占满，所有请求都变慢/超时）。

**一句话**：ROLLBACK 不是修复了什么，而是"正确、干净地收尾一次已经失败的操作"，
是代码写得规范、负责任的体现，不是 bug 的一部分。

### Go/bun 里长什么样

```go
// 简化示意，非原文件逐行摘录
tx.Begin()
count, err := tx.Query("SELECT COUNT(...) FROM communication AS c WHERE ... c.target_app IN (...)")
if err != nil {
    tx.Rollback()   // ← 第一步就失败，必须显式收尾
    return err
}
// ... 第二步查询本该在这里继续 ...
tx.Commit()
```

## 2. health-check ハートビート

### 一句话定义

`/health-check/` 这类接口，是**服务器自己**用来证明"我还活着、还能正常响应"的接口，
**不是**用来检查客户端（浏览器/用户）死没死——方向是反的。

### 谁在调、为什么每 30 秒一次

现实中（这次 floci 本地环境也一样模拟了这个行为），是**负载均衡器（ALB）**
每隔一段时间主动去敲一下每个后端服务的这个接口，问一句"你还活着吗"：

- 连续几次没人应答 / 返回非 2xx → 负载均衡器判定"这个实例挂了"，
  **自动停止把真实用户的请求转发给它**，改发给别的健康实例（如果有多台的话）。
- 这套心跳完全是**基础设施在背景里自动做的事**，跟有没有真实用户在操作
  浏览器无关——这也是为什么日志里没人点任何东西的时候，还是每 30 秒
  冒出一条 `GET /health-check/ 200`，不是针对某次用户操作产生的。

### `/health-check/` 能证明什么、不能证明什么（重要边界）

以 fhir 服务为例：

- **能证明**：进程启动成功、正在 listen 端口。因为路由注册后紧跟着执行的
  DB 初始化（`welbyMicroOrm.NewDatabaseFromConfig`）如果失败，`InitServer`
  会直接返回错误、不会进入 listen 状态——所以 200 间接证明了"启动 + DB 连接成功"。
- **不能证明**：具体某个业务功能（认证、某张表的某个查询）是否正常。
  实测中 `/health-check/` 全部返回 200，但登录后首页依然因为
  `communication.target_app` 这一列不存在而 500——**健康检查通过 ≠ 功能正常**，
  必须靠更深一层的业务级验证（对应 README 里说的"L1 只证明启动，功能确认要看 L2/L3"）。

## 3. 两者在同一条真实日志里长什么样

2026-09-04 实测 floci 时抓到的一条完整链路（`docker logs <fhir容器> --tail 100`）：

```
[bun]  BEGIN                 1.083ms  BEGIN
[bun]  SELECT ... FROM communication AS c ... (c.target_app IN ('all','my_karte'))
       *mysql.MySQLError: Error 1054 (42S22): Unknown column 'c.target_app' in 'where clause'
{"level":"ERROR", "msg":"failed to count information list: Error 1054 ..."}
[bun]  ROLLBACK                474µs  ROLLBACK
{"level":"ERROR", "msg":"failed to get information list: ..."}
{"level":"INFO",  "msg":"usecase error", "err":"failed to get information list: ..."}
... | 500 | POST | /api/v1/mykarte/web/information/from-welby/list | -

... | 200 | GET | /health-check/ | -   ← 每30秒一条，跟上面的500完全独立、互不影响
```

**根因**（跟本文件的 DB 知识点无关，单纯记一下当时的业务背景）：`fhir` 的
`develop` 分支自 2026-08-19（WLB_NEW_KARTE-4769）起依赖 `communication.target_app`
这一列，但 `welby-db-aggregator-platform` 仓库生成 `welby_base` 表结构的
`create_welby_base_db_tables.sql` 最后更新是 2026-08-06，早于该功能合并，
没有这一列——是两个仓库之间的 schema 漂移，不是 floci 或这次 PR 的问题。
