# FHIR リソース種別の基礎知識——マイカルテのコードでどこにあるか

> 业务知识笔记 — 源自 2026-08 TISI（外部監査法人）向け WPDP プラットフォーム審査対応
> （納品物 No.4「主要FHIRリソースと関連テーブル」= e4 項目の調査、新規アサイン時のオンボーディング整理）
> 姊妹篇：[er-diagram-and-fhir-resource-mapping.md](./er-diagram-and-fhir-resource-mapping.md)
> 讲的是 ER 图基础 + FHIR 逻辑映射 vs 物理表定义的区别

## FHIR 是全球标准，不是日本专属

**FHIR**（Fast Healthcare Interoperability Resources）是 [HL7 International](https://www.hl7.org/fhir/)
制定的**国际**医疗数据交换标准，全世界都在用，不是日本专有的。

日本有一个建立在全球版 FHIR 之上的本地化规范 **JP Core**（JP Core Implementation
Guide）——不是另一套独立标准，而是针对日本医疗场景加的约束/扩展（比如姓名的
汉字/片假名表示、日本健康保险番号格式等）。搜过 `welby-fhir-server-aggregator-platform`
仓库的代码和文档，没有出现 "JP Core" 字样——说明 Welby 这边大概率没有严格照搬
JP Core 规范，而是自行设计了一套基于 FHIR resource 概念的映射。跟 TISI 说明
"我们的 FHIR 实现范围"时，这点值得提一下。

## Resource（资源）：FHIR 拆分医疗数据的基本单位

FHIR 把医疗数据拆成一个个类型化的 **Resource**，每种类型代表一类信息。
`welby-fhir-server-aggregator-platform` 仓库里，`infrastructure/repository/fhir/`
目录下的 12 个子文件夹，基本对应 Welby 系统实际用到的 FHIR resource：

```
infrastructure/repository/fhir/
├── patient/            ← 患者基本信息
├── observation/        ← 检查/测量结果（最常用）
├── medication/         ← 药物/处方
├── document_reference/ ← 文档引用（诊断书 PDF 等）
├── questionnaire/      ← 健康问卷
├── organization/       ← 机构信息
├── practitioner/       ← 医疗从业者
├── device/             ← 医疗设备
├── consent/            ← 患者数据使用同意
├── provenance/         ← 数据来源/溯源记录
├── complex_types/      ← 各 resource 共用的字段结构（非独立 resource）
└── data_register_log/  ← Welby 自建的数据登记日志（非 FHIR 标准 resource）
```

每个文件夹里的 `xxx_repository.go` 就是该 resource 具体查询/写入哪些物理表的
代码——**这本身就是"FHIR resource ⇄ 物理表"映射关系，只是以代码形式存在，没有
整理成文档**。这也印证了「没有统一映射文档，信息分散在各 resource 的 repository
代码里」这个判断。

domain 层的定义在 `internal/domain/fhir/`：

- `resource/` —— 各 resource 的 view model（`patient_view_model.go`、
  `observation_search_parameter_view_model.go`、`consent_view_model.go`、
  `document_reference_view_model.go`、`bundle_view_model.go` 等）
- `data_type/` —— 复合数据类型定义

## 逐个 resource 讲什么

**核心业务类（和病人健康数据直接相关）**

| Resource | 是什么 | 举例 |
|---|---|---|
| `Patient`（患者） | 患者基本信息，**整个系统的锚点**——几乎所有其他 resource 都会关联一个 patient，回答"这是谁的数据" | 姓名、性别、生年月日、患者ID |
| `Observation`（观测/检查值） | **最核心、用得最多**，泛指"对患者做的一次测量或检查"，マイカルテ通过设备连携（DDA）拿到的健康数据基本都落在这里 | 血压、血糖、体重、检验值 |
| `Medication`（药物） | 处方药、用药记录 | 开了什么药、用法用量 |
| `DocumentReference`（文档引用） | 指向某个文档文件的**元数据**，不存文件本身 | 诊断书 PDF 的类型/归属描述 |
| `Questionnaire`（问卷） | 健康问卷调查内容 | 生活习惯问卷、症状自评表 |

**机构/人员类**

| Resource | 是什么 | 举例 |
|---|---|---|
| `Organization`（机构） | 医院、诊所、检验机构等组织信息 | 医院名称、所属机构 |
| `Practitioner`（医疗从业者） | 医生、护士等医疗人员信息 | 开处方的医生是谁 |
| `Device`（设备） | 患者使用的医疗器械信息本身——区别于 Observation（设备**测出来的值**），Device 是"这台设备是什么" | 血压计、血糖仪的设备信息 |

**权限/合规类**

| Resource | 是什么 | 举例 |
|---|---|---|
| `Consent`（同意） | 患者的数据使用同意书，医疗数据领域特别重要，涉及隐私合规 | 是否同意把数据共享给某机构/用于某用途 |
| `Provenance`（数据来源/溯源） | 记录"这条数据是谁、什么时候、通过什么方式创建或修改的"——不是健康数据本身，是健康数据的**履历记录**，用于审计追溯 | 一条 Observation 是哪次 DDA 同步写入的 |

**辅助/技术类（Welby 自己的实现分类，不是标准 FHIR resource 名称）**

| 分类 | 是什么 |
|---|---|
| `complex_types`（复合类型） | 多个 resource 共用的"零件"数据结构，比如 `CodeableConcept`（编码+说明的组合，如"高血压"这个诊断既有编码又有文字）、`Quantity`（数值+单位的组合，如"120 mmHg"） |
| `data_register_log`（数据登记日志） | Welby 自己加的"数据什么时候被写入/登记"日志，用于内部追踪，不对外算作 FHIR resource |

## 一句话建立整体图像

一个患者（**Patient**）在某家医院（**Organization**）、由某个医生（**Practitioner**）
看诊，用血压计（**Device**）测出的血压值记录成一条检查结果（**Observation**），
开的处方是 **Medication**，诊断报告文件用 **DocumentReference** 指向，患者是否
同意数据共享记在 **Consent** 里，这一切操作的来龙去脉由 **Provenance** 记录留痕。

要看某个 resource 具体落在哪几张物理表，直接看对应文件夹里的
`xxx_repository.go`——SQL/ORM 查询里会写出具体表名。

## 谁可能懂 FHIR mapping——一条未确认但重要的人脉线索

2026-08-31（一）14:00 团队会议上讨论 e4 任务范围时，有同事提到一个关键信息：

> **高橋さん（Takahashi）可能管理着 FHIR mapping 相关的文档，据猜测放在
> Confluence 上**（原话"高橋さんがもしかしたら管理してるかもしれないんですけれど"
> "クローグにあるかな"——后半句很可能是"コンフルにあるかな"被听岔了）。

**这条线索的性质：不是事实，是二手推测**——说这话的人自己也用了"かもしれない"
（可能），当天高橋さん本人没有在场确认，具体他管理的范围（是所有 FHIR resource，
还是只有某几个，比如提到的 Observation、Organization）也没说清楚。

会议进一步确认了本笔记前面提到的判断：**FHIR mapping 没有统一文档**，是各服务
（マイカルテ、サイログ等）各自独立设计的，知识分散在各服务设计者脑子里 /
Backlog 讨论记录里，当前 TISI 手头的文档完全没有 backlog 链接可追溯来源。

**因此团队的应对策略是**：不笼统接下"给我们完整 FHIR mapping 文档"这种大而化之
的需求（范围会无限膨胀），而是反问对方具体要哪个 FHIR resource（比如
Observation、Organization），范围收窄后，团队成员（包括高橋さん）大概率能针对
具体 resource 给出回答。

**后续行动（尚未做）**：直接找高橋さん本人确认 (a) 他到底负责/了解哪部分 FHIR
resource mapping，(b) Confluence 上是否真的存在这份文档。在确认之前，不能把
"高橋さん管理着 mapping 文档"当作对 TISI 的正式回复依据。
