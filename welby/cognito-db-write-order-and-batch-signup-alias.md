# スタッフ一括登録：Cognito書き込み順序とメールaliasトリック

> 业务知识笔记 — 源自 2026-09-18 review `welby-auth-webserver-aggregator-platform`
> [PR #317](https://github.com/Welby-Inc/welby-auth-webserver-aggregator-platform/pull/317)
> （WLB_NEW_KARTE-5113，org 10105 スタッフ一括登録 runbook）时，逐条对照 `internal/service/staff.go`
> / `user_service.go` / `change_id.go` 源码验证出来的一套认证/批量注册机制。
> 姊妹篇：[db-transaction-rollback-and-healthcheck.md](./db-transaction-rollback-and-healthcheck.md)
> 讲的是「本地事务」本身；这篇讲「本地事务 + 一个不能自动回滚的外部系统（Cognito）」放在一起时怎么排序。

## 1. 业务背景：施設移行 Batch2 是什么

Welby 在把一批诊所/药店（施設）迁移进当前平台，分两阶段：

- **Batch1**（WLB_NEW_KARTE-5090）：批量注册**机构本身**的信息，389 家全部成功。
- **Batch2**（WLB_NEW_KARTE-5113，当前这条笔记的背景）：批量注册这些机构里的**医护人员/管理员账号**，
  用脚本 `_misc/create_organizations/create_bulk_staffs.sh` 读一份 5 列 CSV
  （`name,email,password,role_id,organization_id`），逐行调用 3 个 API：
  `signup` → 等确认邮件 → `confirm-signup` → `setup-staff`。

**org 10105 为什么一开始被排除在外**：业务规则是"一个机构必须刚好有一名管理员（`role_id=3`）"。
7/29 最原始的 CSV 里，org 10105 的 3 名员工全是 `role_id=2`（协理），没人是管理员。准备正式提交
文件的人（白川/山内）发现这个问题后，**在跑脚本之前就把这 3 行整个删掉**，不是"跑了但失败"，
是"压根没交给脚本处理"。直到 9/17 拿到新版 CSV（3 人里 1 人升级成 `role_id=3`）才第一次注册成功。

## 2. 批量注册的核心难题：Cognito 要求"真的收到邮件"

Cognito（认证系统）注册一个新账号，硬性规定必须往邮箱发确认码、必须把确认码交回去才算激活——
不会因为"这是批量脚本"就开后门跳过。而现实是 189 个真实员工的邮箱，操作者一个人根本打不开。

### 用 Gmail 的 `+` 别名，让一个收件箱代收 189 份确认邮件

`create_bulk_staffs.sh:341` `derive_from_email()`：

```
<mailbox>+b2_<org_id>_<line_no>@<domain>
```

比如操作者邮箱是 `yusuke.saito@welby.jp`，org 10105 的第 2/3/4 行会分别生成：

| CSV 行号 | 临时登录地址（喂给 Cognito 的字符串） | 实际落地的收件箱 |
|---|---|---|
| 第2行 | `yusuke.saito+b2_10105_2@welby.jp` | 还是 `yusuke.saito@welby.jp` |
| 第3行 | `yusuke.saito+b2_10105_3@welby.jp` | 还是 `yusuke.saito@welby.jp` |
| 第4行 | `yusuke.saito+b2_10105_4@welby.jp` | 还是 `yusuke.saito@welby.jp` |

这里其实是**两件不同的事叠在一起**，容易混：

1. **"189个不一样的地址=189个不同的人"——跟 Gmail 无关，是 Cognito 自己的行为。**
   Cognito 很死板：只要字符串不一样，就当成不同的人，建出不同的账号（不同的 `sub`）。
   它不知道、也不关心这些地址背后是不是同一个真实邮箱。
2. **"这189封发给假地址的确认邮件不会退信、真的有人能收到"——这才是 Gmail 的 `+` 别名规则。**
   Gmail 的邮件服务器认得 `+` 后面的内容可以忽略，自动把信投递到 `+` 前面那个真实存在的邮箱，
   同时把完整地址保留在邮件的收件人信息里，供后续区分。

脚本读确认码时（`create_bulk_staffs.sh:530-532`），不是去收件箱里随便翻，而是精确搜索
`to:${from_email} AND is:unread`——即使收件箱里同时躺着 189 封长得很像的确认邮件，也能精确
挑出属于"这一行"的那一封。

**两个条件缺一个都不行**：Cognito 不认"不同字符串=不同人"，这套办法从根上走不通；
Gmail 不认 `+` 号规则，189 封确认信会全部退信，流程照样卡死。

### `+` 规则是 Gmail 专属，还是所有邮箱都有

不是互联网邮件的统一标准，是**各家邮箱服务商自己的实现选择**（正式名字叫
sub-addressing / plus addressing，邮件地址本地部分允许 `+` 字符是 RFC 允许的，
但"收到 `+` 要不要忽略后面的内容"完全是各厂商自己决定）：

| 邮箱服务商 | 支持情况 |
|---|---|
| **个人 Gmail（`@gmail.com`）** | 支持，而且是**最早、最主要的使用场景**——Google Workspace（企业版，`welby.jp` 用的就是这个）跑的是同一套底层引擎，继承了同样的能力，**没有管理员开关能单独关掉** |
| **Outlook.com / Microsoft 365** | 现在大部分个人版也支持了，但**企业版很多租户默认关闭**，需要管理员手动开——跟 Gmail 正好反过来（Gmail 企业版默认开、微软企业版默认可能关） |
| **iCloud Mail（Apple）** | 不支持 `+`。苹果给了一个更彻底的替代方案"隐藏我的邮箱"（Hide My Email）：生成一个跟真实地址**完全不相关的随机地址**转发，比 `+` 标签更彻底——`+` 标签的问题是任何人看到 `你+淘宝@gmail.com` 都能一眼猜出真实地址是 `你@gmail.com`，直接掐掉标签就能拿到，并不是真正的隐私保护，只是方便自己筛选邮件 |
| **Yahoo Mail** | 历史上用的是另一套语法不同的"一次性地址"功能，跟 `+` 写法不一样 |
| **企业自建邮件服务器**（Postfix/Exchange 等） | 完全看管理员有没有配置 |

**Gmail 为什么会有这个功能**：`+` 这个玩法不是 Gmail 发明的，是 1990 年代 Sendmail/Postfix
这些 Unix 邮件服务器就有的老传统，Gmail 后端团队顺手继承过来——实现成本低（投递时看到 `+`
忽略后面部分即可），换来一个用户能自己用来筛选/追踪邮件来源的实用功能，对 Gmail 是低成本高收益。

**其他家为什么没跟进**：不是技术做不到，是取舍不同——iCloud 选了"更彻底但成本更高"的隐私方案；
企业邮箱系统（微软 365/Exchange）更谨慎，是因为很多老系统、内部工具几十年来假设"邮箱本地部分
就是纯用户名，没有特殊字符"，默认开启可能悄悄破坏老系统的路由逻辑；也存在被滥用的场景——
如果公司内部某个校验逻辑写得不严谨（比如判断"邮箱字符串里包不包含某个关键词"而不是"跟白名单
地址精确相等"），域名一旦支持 `+`，就有可能被人拼出一个标签里塞了关键词的地址去蒙混过关。
（这只是说明性例子，不是某个具体已知漏洞编号，用来解释企业邮箱系统为什么对默认开启这件事谨慎。）

**CSV 里的 `password` 是什么**：不是占位符/测试密码，是这个员工**最终真实能登录用的密码**，
批量脚本原样把这一列传给 `signup`。这也是为什么整份 runbook 反复强调"密码明文、不能进 git、
跑完要删 CSV"——泄露这份 CSV 等于直接公开了一批真人账号的密码。

## 3. Cognito 与本地 DB 的写入顺序：为什么"先 DB 后 Cognito"是对的，不是 bug

`internal/service/staff.go` 的 `SetupStaff`（"把临时登录地址换成员工真实邮箱"这最后一步）：

```go
// Step 6: Email swap in DB (must occur BEFORE commit)
//   welby_auth.user / welby_base.contact / userapi.users 三处改名

// Step 7: Commit all 3 transactions
txForWelbyAuth.Commit()
txForWelbyBase.Commit()
txForUserAPI.Commit()

// Step 8: Cognito email update (AFTER commit)
s.CognitoClient.UpdateCognitoEmail(*req.FromEmail, *req.ToEmail)
```

即：**先把 3 个数据库的记录全部写完、提交保存，最后才去改 Cognito 那边的登录邮箱**。

Review 时曾讨论过"要不要把顺序倒过来（先调 Cognito、成功了才写 DB，失败就不写，更快失败）"，
结论是**顺序不该换**，理由：

1. **DB 事务失败可以自动回滚**（代码里 `defer func(){ if !committed { rollback() } }()` 保证），
   Cognito 的 `AdminUpdateUserAttributes` **没有任何自动撤销手段**。"改不了的那个动作"应该放最后，
   这样万一它失败，前面的东西不是"已经提交的完整数据"就是"安全回滚的空"，不会出现悬空状态。
2. **倒过来不会消除问题，只会把"谁先落地"换个边，而且换到更糟的一边**：Cognito 先成功、DB 后失败，
   会变成"这个人 Cognito 上能登录，但 DB 里完全没有对应的员工档案"——比现在更难发现、更难诊断。
3. **倒过来会破坏现有的恢复手段**：现在的恢复流程依赖"用 `from_email` 还能在 Cognito 里查到这个人"。
   如果 Cognito 先改了，`from_email` 这个别名就已经失效，下次想用它去找人、去补救——找不到了。

## 4. "昵称" vs "身份证号"：`from_email`/`to_email` 与 `CognitoUserID`（sub）

| | 是什么 | 会不会变 | 谁写入的 |
|---|---|---|---|
| **`from_email` / `to_email`**（昵称） | Cognito 里用来"找人"的邮箱别名 | 会变（注册时是`from_email`，最后swap成`to_email`） | 人工填在 CSV 里，可能打错字/混入怪字符 |
| **`CognitoUserID` / sub**（身份证号） | AWS Cognito 自动生成的固定用户 ID | 永远不变 | AWS 系统自己生成，人不经手 |

`internal/service/staff.go:199` 的 `UpdateMfaEnabledBySub` 已经在用不变的 `user.CognitoUserID`
定位用户；但 `UpdateCognitoEmail`（`infrastructure/aws/cognito_client.go:143`）目前还是用会变的
`from_email` 当 `Username` 参数去定位。**建议**：把 `UpdateCognitoEmail` 也换成用
`user.CognitoUserID` 定位——不需要动第 3 节说的顺序，纯粹是换一个更稳定的"怎么找到这个人"的方式，
以后无论邮箱别名处于什么状态，都能确定找到同一个人。（已作为技术债思路写进 PR #317 的 review 评论里。）

### 准确的因果关系：不是"昵称变成身份证"，是"身份证从没变过，只是昵称被换了一次"

容易说反的一点：**不是"昵称最后变成了身份证"**，昵称和身份证是**两个一直同时存在、各自独立**
的属性，不存在谁变成谁。准确顺序是：

1. `signup` 用 `from_email` 第一次注册的那一刻，Cognito 给这个账号发了一张**身份证**（`sub`）——
   这件事**只发生一次**，之后永远不变。
2. 发身份证的**同时**，还给账号挂了一个"昵称"属性，当时的值是 `from_email`。身份证和昵称从这一
   刻起就是两个独立存在的属性。
3. `UpdateCognitoEmail`（"换邮箱"这一步）做的事，是**找到那张身份证对应的人，把他身上挂着的
   昵称属性从 `from_email` 改成 `to_email`**——身份证本身完全没动，只是昵称这一个属性被换了。

所以 Cognito 把改名前、改名后当成同一个人，**原因是身份证从头到尾没换过**，不是"昵称升级/
转换成了身份证"。类比：一个人办了身份证后改自己的微信昵称，昵称怎么改身份证号都不变，
系统认得是同一个人，是因为身份证这张凭证没换，不是昵称变成了身份证。

## 5. org 10162 事故复盘：为什么"哪个 API 报错"不能用来判断数据写到哪一步

**时间线**（WLB_NEW_KARTE-5113）：

1. 2026-08-13：Batch2 首次本番执行 189 行，162 成功，27 失败。
2. 2026-08-24：修正后重跑 22 行，**21 成功，唯独 org 10162 失败**。
3. **根因**：org 10162 那行的 `email`（CSV 第2列）本地部分里，混入一个肉眼几乎看不出区别的
   半角字符 `ｰ`（U+FF70，半角カタカナ長音記号），不是 ASCII 的 `-`（U+002D）。格式校验能通过，
   但 `setup-staff` 处理到这个字符时返回 HTTP 500。
4. **当时（8/24）的错误推断**：这一批 22 行分配的 `welby_id` 是连续的 564666–564687（22 个号），
   但成功只有 21 行——中间缺了 `564679`，正好对应 org 10162 这一行。团队据此推断
   "practitioner 档案没建成，只是号被占了"。
5. **今天（9/18）的更正**：结合第 3 节的 Step 6→7→8 顺序重新读代码才发现，**"`setup-staff` 500"
   这个观测，根本不能用来判断数据写到哪一步**——因为 3 个数据库是在 Step 7 就已经全部提交，
   Step 8 的 Cognito 更新失败才是 500 的真正原因。也就是说 org 10162 很可能"员工档案已经完整
   建好，只是 Cognito 那边的登录邮箱还停在临时昵称"，跟 8/24 当时猜的完全不是一回事。
   **`welby_id` 编号连续与否，也不能当作"没有部分写入"的证明**——Cognito 更新失败的行，
   编号照样是连续分配的。
6. **现状（未解决）**：还没等到修正版 CSV（把那个特殊字符换成正常的 `-`）。即使拿到了，也
   **不能直接盲目重跑**，必须先实查 Cognito（用户存不存在、`CONFIRMED` 没有、邮箱属性是不是
   已经变成 `to_email`）和 3 个数据库的真实状态，再决定该走"整行重跑" / "单独补 `confirm-signup`"
   / "单独补 `setup-staff`" / "只补 Cognito 邮箱" 里的哪一条，否则可能造出重复档案或账号错乱。

**一句话总结这次学到的教训**：跨系统（本地 DB + 远程 Cognito）没有真正的分布式事务时，
**API 返回的错误码只能告诉你"这次调用没成功"，不能告诉你"之前提交过的东西有没有生效"**——
必须回去看实际执行顺序（哪一步在提交之前、哪一步在提交之后）才能判断真实状态，
猜测式的"从错误码反推数据状态"是这整起事故里反复踩坑的根源。
