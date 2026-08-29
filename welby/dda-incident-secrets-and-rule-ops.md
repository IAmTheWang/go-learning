# DDA Incident: Secrets, Rule Enable/Disable, and Backfill Ops

> 业务知识笔记 — 源自 2026-08 DDA（Device Data Aggregator）生产故障排查（续篇）
> 姊妹篇：[fitbit-cursor-vs-sliding-window.md](./fitbit-cursor-vs-sliding-window.md) 讲的是
> Fitbit 数据同步机制本身；这篇记录故障处理过程中的运维/概念细节。

## 故障时间线速览

1. **2026-08-25 22:57**：生产部署，`codedeploy.sh` 的 secret 注入逻辑只覆盖了
   `internal-api`/`access-api` 两个 family，漏了 vendor batch 的 7 个 family。
   Fitbit / OMRON connect / アークレイ / MeDaCa 的认证 secrets 变成空值，调用厂商
   API 全部返回 401。
2. **故障期间（~3 天）**：这几个 vendor batch 的调度 rule 仍是 ENABLED，按原计划
   反复重试、反复 401，累积出大量失败请求。
3. **2026-08-28**：**タニタ 主动发邮件投诉**"你们的生产 IP 在大量访问我们"——
   这正是故障被发现的契机（厂商先报警，不是 Welby 自己先发现）。
4. **排查**：中村さん牵头查日志，余さん定位到根因——同一次部署里
   `internal-api`(secrets=12)/`access-api`(secrets=12) 正常，只有 7 个 vendor
   batch family 是 `secrets=0`。
5. **止血**：
   - タニタ：手动临时注入 secret，先复活并立刻回填了 3 天缺口。
   - Fitbit/OMRON/アークレイ/MeDaCa：持续 401 无法立刻手动修，先把调度 **rule
     手动 DISABLE**，让它们别再打厂商 API。
6. **根治**：PR #105（修复 `codedeploy.sh` 密钥注入范围）+ PR #107（Fitbit 专用
   回填脚本，因为 Fitbit 没有游标，重新 enable 不会自动补数据）。staging 验证后
   走审批，部署到生产。
7. **恢复**：部署后确认 6 个 vendor batch 的 task-def 都拿到 `secrets=12`，再把
   之前 DISABLE 的 4 个 rule 手动 ENABLE 回去，首轮运行 401 全部为 0。
8. **收尾**：Fitbit 的历史数据回填（WLB_NEW_KARTE-5350）先跑 DRYRUN 确认范围，
   发现 601 名目标患者中 252 人（42%）的个人 Fitbit 授权 token 已过期
   （`invalid_grant`），这是跟本次故障无关的另一个问题，需要另外处理。

## `secrets=N` 是什么意思

`secrets=12` 里的数字，指的是 **AWS ECS 任务定义（task definition）里挂载了多少条
密钥**，不是什么特殊编号。`codedeploy.sh` 部署时会把密钥从 Secrets Manager 取出，
写进 task definition 的 `secrets` 数组，每一项是一条"环境变量名 → 密钥 ARN"的映射：

```json
"secrets": [
  { "name": "FITBIT_CLIENT_SECRET", "valueFrom": "arn:aws:secretsmanager:...:fitbit-secret" },
  { "name": "OMRON_API_KEY",        "valueFrom": "arn:aws:secretsmanager:...:omron-key" },
  { "name": "DB_PASSWORD",          "valueFrom": "arn:aws:secretsmanager:...:db-pass" }
]
```

`secrets=12` 就是这个数组里有 12 条映射（可能包含各厂商 API key、数据库密码、内部
服务凭证等）；`secrets=0` 就是数组是空的，容器启动后所有需要密钥的地方全是空值。

这个数字在排查和验证时都被用作最直接的"体检指标"：不用去猜密钥内容对不对，只看
数字对不对就知道注入逻辑生效没有。

- 故障时：`internal-api`/`access-api` 是 `secrets=12`（正常），7 个 vendor batch
  是 `secrets=0`（异常）——这个数字差异就是定位根因的直接证据。
- 修复后：Fitbit/OMRON/アークレイ/MeDaCa/タニタ/A&D 全部变成 `secrets=12`，证明
  PR #105 生效。
- マイナポータル 是 `secrets=3`——数字不同是正常的，它走 terraform 单独管理，
  本来需要的密钥种类就比其他 vendor batch 少，不是 bug。

## 两种容易混淆的"token/密钥"

| | 是什么 | 归属 | 这次故障里的角色 |
|---|---|---|---|
| **secrets**（app 级密钥） | Welby 系统证明"我是 Welby"的 API Key/Client Secret | 属于整个应用，跟具体患者无关 | **这次故障的根因**——`codedeploy.sh` 漏配置导致变空 |
| **OAuth token**（用户级授权） | 每个患者自己在 Fitbit 上授权后拿到的 access/refresh token | 属于每个患者个人 | **跟这次故障无关的另一问题**——601 人里 252 人这个 token 早已过期失效（`invalid_grant`），到现在还没解决，需要患者重新做 Fitbit 连携授权 |

## rule 的 enable/disable：正确的因果顺序

容易搞反的一个点：**disable 不是"为了强行读到数据"，而是"为了让它别再读了"。**

1. **故障发生**：secrets 变空，但 rule 仍是 ENABLED → 批处理按原计划继续跑，
   持续 401，持续制造失败请求（这正是被タニタ 当成"大量访问"投诉的原因）。
2. **临时止血**：手动把 rule 从 ENABLED 改成 **DISABLED**——目的是让"没密钥还硬
   敲"的循环立刻停下来，纯粹止损，跟"读数据"无关。
3. **根治修复**：PR #105 让 `codedeploy.sh` 正确注入 secrets，重新部署后
   task-def 拿到 `secrets=12`。
4. **恢复运行**：secrets 已经修好，这时候才把 rule 从 DISABLED 改回
   **ENABLED**——因为现在再跑就不会 401 了。重新 enable 后首轮运行确认 401=0。

顺序是：**先止血（disable）→ 修根因（改代码部署）→ 恢复（enable）**，而不是反过来
"为了强行拿数据所以 disable"。

## 为什么要做"限流保护的加固"

这是故障复盘后新增的一项后续加固工单（跟本次故障根因无关，是发现的次生风险）。

**是怎么被发现的**：secrets 变空后，批处理任务不会自己停下来，而是按原调度频率
（约每 15 分钟一次）持续敲厂商 API，每次都被 401 拒绝，再等下一轮继续敲。乘以几百
个患者、乘以将近 3 天、乘以每 15 分钟一次的调用频率，失败重试的总请求量堆积得相当
可观——这正是タニタ 监控到"大量访问"、误判成攻击或爬虫，主动报警的根本原因。

**要解决的问题**：当前系统**没有任何"认证持续失败就自动降频/熔断"的机制**。加固
的方向是：连续 N 次 401 之后，批处理任务应该自动降低重试频率或直接暂停，而不是
一直按正常节奏死磕厂商 API。这样以后即使再出现类似的密钥配置错误，也不会在厂商
那边制造出"看起来像攻击"的流量洪峰，能更早自证清白，也保护对方服务器不被误伤。

## 为什么生产环境的一次性回填脚本要先 DRYRUN，不能直接 write

Fitbit 回填（PR #107）风险和影响面都不小，而且是专门为这次故障新写的一次性脚本，
没有经过大量线上验证：

1. **写的是真实患者的医疗数据**：这次要往 601 名患者的マイカルテ（健康档案）里
   补约 5 万条 Observation 记录。日期范围算错、患者过滤条件写错，写进去的可能是
   重复数据或覆盖了不该覆盖的记录，脏数据清理远比"跑之前先检查"麻烦。
2. **提前发现风险，避免带着错误假设写库**：DRYRUN 只做"拉取 + 解析"，跳过
   "写 lifelogs/写 FHIR"这一步，相当于空跑一遍看范围、数量级对不对。这次 DRYRUN
   真的挖出了本节前面提到的 252 人 token 失效问题——如果没有先 DRYRUN，直接跑
   正式回填，很可能跑到一半才发现这批人全部失败，还得现场排查是脚本问题还是
   患者授权问题。
3. **一次性脚本缺乏大量线上验证**：staging 阶段只用少量测试患者验证过（且这些
   患者都因为 token 失效没能走到最后一步），到生产环境第一次对全量数据跑，先
   DRYRUN 是最基本的"没底的东西先空转一遍"的稳妥做法。

一句话：DRYRUN 成本几乎为零（不写库），但能提前暴露范围错误、数量异常、意料之外
的失败，比起直接写库出问题再回滚，风险收益完全不对等。
