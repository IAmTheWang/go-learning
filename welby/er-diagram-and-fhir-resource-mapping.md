# ER図・Entity-Relationship と FHIR リソースマッピングの基礎知識

> 业务知识笔记 — 源自 2026-08 TISI（外部監査法人）向け WPDP プラットフォーム審査対応
> （納品物 No.4「主要FHIRリソースと関連テーブル」= e4 項目の調査）
> 姊妹篇：[device-integration-service-pattern-and-git-forensics.md](./device-integration-service-pattern-and-git-forensics.md)
> 讲的是機器連携サービスの設計パターンと git 調査手法

## Entity-Relationship（ER）到底是什么——用前端类比理解

**Entity（实体）= 一张数据库表，约等于前端写的一个 `interface`/`type`。**

```ts
interface Patient {
  patientId: string;   // PK 主键
  addressId: number;   // FK 外键——指向 Address 表的某一行
  contactId: number;   // FK 外键
}
interface Address {
  addressId: number;   // PK
}
```

`Patient.addressId` 不是自己存地址信息，而是存一个"指针"（外键），指向 `Address`
表里的某一行。这跟前端写 `patient.addressId` 再拿这个 id 去另一个数据源查 `Address`
对象是一回事——数据库层面把这种"谁引用谁"的关系正式记录下来，并强制约束（不能往
`addressId` 里塞一个 `Address` 表里不存在的值）。

**Relationship（关系）就是这些 interface 之间怎么互相引用**，ER 图就是把"表 +
引用关系"画成图，而不是让人去翻几百条 `CREATE TABLE` 语句。

Mermaid 的 `erDiagram` 语法里，这样一行：

```
patient }o--|| address : "address_id (ext)"
```

翻译成人话："多个 patient 可以对应 1 个 address"（`}o--||` 是基数符号，多对一）。

## Mermaid：用纯文本写图的工具

不需要画图软件（Visio/draw.io）手动拖拽连线，而是写一段像代码一样的文本：

```
erDiagram
  patient {
    varchar patient_id PK
  }
  address {
    bigint address_id PK
  }
```

GitHub / VS Code 看到 ` ```mermaid ` 代码块会自动渲染成图形。好处：**能跟着代码一起
用 git 管理版本**——表结构改了，图的定义也在同一次 commit 里改，历史可追溯，不用
额外维护一份画图软件的文件。Mermaid 不止能画 ER 图，还能画流程图、时序图（sequence
diagram，画的是"一次操作里几个角色按时间顺序谁给谁发了什么消息"，跟 ER 图的静态
结构完全是两回事）、甘特图、类图等。

## welby_base 的 ER 文档集：docs-as-code 的实践

项目里 `welby-fhir-server-aggregator-platform` 仓库的 `_deploy/designDocs/er/`
目录，是**一套自动从 `information_schema` 生成的 ER 文档**：

- `README.md` / `DRIFT.md`（记录已知的文档-实际环境差异）/ `welby_base.dbml`
- `diagrams/*.md`：34 个文件，18 个业务主题域（Patient/Person、Care Plan & Goal、
  Clinical Results、Medication、Consent & Contract、Document、Device、
  Audit & Provenance、Migration Staging 等），每个主题域一份 plain 版 + 一份
  `__history` 版
- `scripts/regenerate.sh`：一键从当前 DDL 重新生成全部 Mermaid/DBML

**docs-as-code 的核心价值**：文档不是人工画的，是脚本从真实 schema 抽取生成的。
这意味着"文档过时"这件事，理论上只需要重跑一次脚本就能自动修复——不需要人工逐张表
去核对。截至 2026-08-01（齋藤さん，WLB_NEW_KARTE-5008）生成时是 562 张表 /
771 外键；到 2026-08-31 实际 DDL 已经变成 564 张表 / 774 外键，6 张表被移除
（`generic_note`/`intervention`/`organization_dcf_facility_code` 等）、2 张表
新增但未反映到文档（`patient_tsuno_link`、`external_service_patient_link`）——
这就是"脚本能自动生成但没人重新跑"导致的典型 docs-as-code 漂移案例。

## FHIR 资源逻辑映射 vs 物理表定义——两者的本质区别

**物理表定义书**：单纯列出"这张表叫什么、有哪些列、列的类型是什么"——是数据库层面
的事实，跟业务语义无关。

**FHIR 资源逻辑映射表**：把每一张物理表，归类到 [FHIR](https://www.hl7.org/fhir/)
（医疗数据交换的国际标准）定义的资源类型上，比如 `patient` 表对应 `Patient` 资源、
`observation` 表对应 `Observation` 资源、`medication_request` 对应
`MedicationRequest`。这是一层**业务语义分类**，没法从 `information_schema` 自动
抽取——机器只能告诉你"这张表叫 `observation`"，但没法自动判断"这张表在 FHIR 语义
里算 `Observation` 还是 `DiagnosticReport` 的一部分"，这需要懂 FHIR 规范 + 懂业务
含义的人来做判断。

这也是为什么外部审计方（TISI）明确指出："物理表定义已经很详细了，但看不出跟
『主要FHIRリソース』的对应关系"——这两种文档解决的是完全不同层次的问题，交了物理
表定义并不能替代 FHIR 资源映射。

实际调查中见过的一份 425 行草稿映射表，逐行标注了置信度：

| 确认状态 | 数量 | 含义 |
|---|---|---|
| 確定 | 322 | 有把握的分类 |
| 推定 | 68 | 大致判断，未完全确认 |
| 要確認 | 35 | 需要进一步人工核实 |

这种"分类 + 自评置信度"的三段式结构，如果是人工做的说明是审慎的专业判断；但如果是
喂给 AI 工具自动生成的，"確定"也只是 AI 自己觉得有把握，不代表真的对——**来源不明
的分类结果，置信度标注本身不能替代人工复核**，尤其是要拿去交给外部审计方的场合。

## 物理表 vs FHIR resource——一个仓库比喻

初次接触这两个概念很容易搞混，用仓库货架来类比会直观很多：

- **物理表（physical table）= 你们仓库里真实存在的货架**，怎么摆、分几层、用什么
  材料，完全是工程师自己根据业务需求、数据库性能考量自由设计的（对应：`observation`
  表有 `id`、`patient_id`、`value`、`unit`、`measured_at` 这些列，全是自己公司决定的）。
- **FHIR resource = 国际组织（HL7）发布的通用货物分类标准**，不是某张具体的表，
  而是一份规范文档，规定"凡是『对病人做的一次检查/测量』这种信息，不管哪个国家
  哪个系统，都应该归为『Observation』这一类"。它不认识你数据库里那张表叫什么
  名字，只关心"这堆信息在概念上属于哪一类"。
- **为什么需要人工判断、人工"规定"对应关系**——因为国际标准不会自动知道你们的
  `observation` 表该算 Observation（哪怕碰巧同名，也可能叫别的名字如 `exam_result`）。
  需要一个懂业务 + 懂 FHIR 规范的人，把物理表的定义一张张过一遍，判断"这张表按
  FHIR 标准该归为哪一类"，再写下来——**这个写下来的对应关系就是 FHIR mapping
  文档**。

**关键纠偏**：不是"物理表必须在结构上符合 FHIR 标准"，而是"物理表可以自由设计，
只需要能在需要时被翻译成 FHIR 格式"。仓库货架怎么摆是自己的事，FHIR 只要求
"外部来查货时，你能说清楚这批货换算成国际分类标准是哪一类"。仓库里
`xxx_repository.go` 查物理表、拼出 FHIR resource 对象返回的过程，就是这个"翻译"
实际发生的地方。

## 分类粒度：按表分类，还是按字段分类？

继续用仓库比喻：分类的**主体单位是"整个货架"，不是"货架上的某一格"**，但格子这层
也牵涉进来，只是发生在"给货架贴大类标签"**之后**：

1. **先给货架整体贴一个大类标签**：仓库管理员先判断"这一整个货架，按国际分类标准
   该归为哪一类"——比如 `observation` 这个货架，整体被判定为"Observation 类货架"。
   这一步是**货架级**的，一次判断，管一整个货架。
2. **再给货架里每一格贴小标签**：确定这个货架属于"Observation 类"之后，才回过头
   来看货架上摆的每一格具体货品（表的每一列），对应国际标准里"Observation 类"
   规定要有的具体细项——比如"这一格是编号"→标准要求的"subject"，"这一格是数值"
   →标准要求的"value"，"这一格是时间"→标准要求的"effectiveDateTime"。格子本身
   不会脱离它所在的货架，单独被搬去贴别的大类标签。

**两个例外，用具体例子讲清楚**：

**例外1：一张表里混装了好几种不同的东西，不能整张表算一类**

假设有一张表长这样：

| id | type | patient_id | content | created_at |
|---|---|---|---|---|
| 1 | vital_sign | 001 | 血压120 | 2026-08-01 |
| 2 | consent | 001 | 同意共享数据 | 2026-08-02 |
| 3 | vital_sign | 002 | 体重65kg | 2026-08-03 |

这张表**不能整张表统一说"它是 Observation"或"它是 Consent"**——因为第 1、3 行是
检查数据（该算 `Observation`），第 2 行是同意书（该算 `Consent`）。这种情况下，
要看每一行的 `type` 这一列的值，才能知道**这一行**该归到哪个 FHIR 分类，不是整张
表一次性归类。判断的单位从"一整张表"缩小到了"表里符合某个 `type` 值的一批行"。

（现实中这种设计通常出现在比较老、比较图省事的表结构里，一张表塞很多种不同类型
的记录，靠 `type` 字段区分——不是好的设计习惯，但确实存在。）

**例外2：一个 FHIR 分类需要凑好几张表的内容才够**

假设 FHIR 标准规定："`Patient` 这个分类必须包含姓名 + 地址 + 电话"。但数据库里
可能是分开存的：

- `patient` 表：只有姓名
- `patient_address` 表：只有地址
- `patient_contact` 表：只有电话

要拼出一个"完整、够格"的 `Patient` resource，得**同时查这三张表**，把姓名、地址、
电话拼在一起，才够 FHIR 标准要求的内容——对应到代码里，就是 join 这几张表的字段
才能构造出一个完整的 FHIR `Patient` resource 对象。这种情况下，**一个 FHIR 分类
不是对应一张表，而是对应好几张表拼起来的结果**。

一句话总结两个例外：
- 例外1：一张表里混了好几类东西，要看"行"才能分类（不是整张表一类）
- 例外2：一个 FHIR 分类的内容分散在好几张表里，要拼起来（不是一张表能装完）

Welby 目前仓库里 `xxx_repository.go` 的实现基本是第一种模式（表→resource 整体
对应），这也是为什么本篇前面说"FHIR 资源逻辑映射是表级别的语义分类"是主要模式。
