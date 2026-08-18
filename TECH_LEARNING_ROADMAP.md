# Tech Learning Roadmap（AWS / Docker / Go / PostgreSQL / 英语）

> 这份文件是从 `~/personal/resumes` 项目里沉淀的求职策略文档中抽取出来的技术学习方向指引，
> 同步放在 `go-learning/` 和 `docker-learning/` 这两个技术练习仓库里，作为"为什么学这个、
> 优先级是什么"的背景说明。完整的市场调研 / 薪资对比 / 理由细节见 resumes 项目的
> `Backend_Language_Strategy_2026-08.md`、`LinkedIn_Job_Leads_2026-08.md`、
> `AU_NZ_Visa_Job_Notes_2026-08.md` 与 `STAR_Stories_EN.md`。

## 背景

前端工程师（Vue 3 / React + TypeScript），目标市场：日本（现职，1-2年内计划转出）+ 中国 +
新西兰 + 澳洲。目前在东京工作两年多，日语已通过N1；因日本2026年对外国人永居政策收紧+
工资相对偏低，计划1-2年内转往新西兰或澳洲。当前AWS和后端都是薄弱项，策略是优先补强
AWS/云原生，不主动深入学习传统后端语言（尤其是Java）；英语（TOEIC 710，工作中使用比例
几乎为0）是新增的并行track，详见下方"英语学习"一节。

## 综合优先级：Docker > AWS SAA > (Terraform深化+EKS) > BFF(Go) > Python > Java

| 排名 | 项目 | 定位 |
|---|---|---|
| 1 | **Docker** | 最高优先级，地基 —— `docker-learning/` 对应这一步，不依赖任何 AWS 知识，可以最先独立学 |
| 2 | **AWS SAA 认证** | 第二优先级，建立 HR 筛选门槛，同时是 EKS 的前置知识（VPC/Subnet/IAM OIDC） |
| 3 | **Go（BFF旗舰项目的载体语言）** | `go-learning/` 对应这一步。**2026-08-16更新：Go不再是独立于BFF之外的加分项，而是BFF旗舰项目本身的实现语言**——学Go即是在做这个项目，定位详见下方"Go 的定位"一节 |
| 4 | Python | 现状够用（Flask AI 辅助维护经验打底），不额外投入 |
| 5 | Java | 不主动学，工作中遇到靠 AI 工具读懂即可 |

## 详细学习顺序（2026-08-16 第二轮优化版）

> 相比最早的"六步计划"，第一轮调整做了三处：解耦SAA与EKS、砍掉AWS CDK、把BFF具体化成
> 旗舰项目。第二轮（本版）在此基础上追加三处：**BFF旗舰项目的语言从Hono/TypeScript改为
> Go**、**Phase 3动手写代码前插入SQL基础复盘**、**新增英语学习并行track**。

| Phase | 核心内容 | 交付产物 | 周期 |
|---|---|---|---|
| **1. 容器与云基础** | Docker/docker-compose 基础 + AWS SAA 刷题体系 | 考出 **AWS SAA 认证** | 1–1.5 个月 |
| **2. IaC 与云原生** | 深化 Terraform + EKS 部署（IRSA/Ingress/Deployment/Service），此时已有 SAA 打下的 VPC/IAM/OIDC 基础，理解 EKS 的网络和权限模型会顺畅很多 | Terraform 自动化拉起 AWS 资源 + 一个 EKS 部署项目 | 1 个月 |
| **3. BFF 实战闭环** | **SQL基础复盘（1-2周，见下）→ Go（Gin 或 Echo）+ sqlc + PostgreSQL**，配合 Docker 做全链路自动化交付 | 完整全栈 Demo，**带 CI/CD 自动部署到 EKS**——这是跟纯 Dev/纯 DevOps 区分开的杀手锏 | 1.5 个月 |
| **4.（可选）AWS DevOps 认证** | 前 3 个 Phase 顺利完成后再考虑，不强制排期 | AWS Certified DevOps Engineer 认证，针对 17-25K 档位那类明确要 CI/CD+运维深度的 JD | 视情况 |

**AWS CDK**：不列入学习目标。处理方式和 Java 一致——不主动投入学习时间，真遇到明确要 CDK
的 JD，靠已有的 IaC 概念（Terraform）+ AI 辅助现学现用即可。

**Java**：不列入学习计划，工作中遇到靠 AI 工具（Claude/Cursor）辅助读懂即可。

## Go 的定位：云原生全栈，而非"Go后端工程师"（2026-08-16）

- **简历/自我定位**：不写"Go后端工程师"，而是"资深前端/全栈工程师
  （React/TS + Go/AWS/DevOps）"或"云原生全栈工程师"。
- **为什么不定位成Go后端专精**：澳新市场主流后端栈其实是C#(.NET Core)、Java、Node.js/TS，
  纯Go岗位偏小众，竞争者多是本地资深后端专精工程师——用"Go后端工程师"身份去正面竞争没有
  优势。Go在Docker/Kubernetes/Terraform这条云原生工具链上是同源语言（这些工具本身都是Go
  写的），这才是Go对这份简历真正的价值：帮助看懂工具内部逻辑和报错信息、以及作为BFF项目的
  实现语言，而不是去竞争"独立Go后端架构师"这类岗位。
- **BFF旗舰项目（Phase 3）的评价标准**：不是"展示Go语法掌握程度"，而是"展示用
  Go+Docker+AWS EKS搭建云原生BFF架构、并完成自动化CI/CD的综合全栈工程交付能力"——技术
  选型是Go，但要证明的核心能力是云原生全栈交付，不是单一语言的语法熟练度。

## Phase 3 前置热身：SQL 基础复盘（1-2周）

回应"MySQL长期不用，基础很弱"这个已知短板——澳新中高级技术面试对数据建模/数据库底层
机制的考察力度很大，只会用sqlc生成代码但答不上底层原理会在面试中翻车。**只啃4个硬骨头，
不漫无目的看教程**（sqlc工具链已经帮你处理好大多数代码生成，语法细节不是重点）：

1. **B-Tree索引失效场景**——对索引列做函数计算、隐式类型转换、前导`%模糊查询`、以及
   违反最左匹配原则导致索引断档，这几种写法怎么让索引失效。
2. **慢查询EXPLAIN分析**——在PostgreSQL里关注`Seq Scan`→`Index Scan`的提升空间，以及
   `Sort`/`Hash Join`/`Nested Loop`等节点是否触发了额外的排序或临时结构。
3. **N+1查询问题**——把循环里逐条查询的写法，换成JOIN或sqlc生成的批量`IN`语句，对比
   查询次数和耗时的差异。
4. **四种事务隔离级别的并发异常**——开两个终端窗口跑`psql`，手动交替执行两个事务的SQL
   语句，直观观察脏读、不可重复读、幻读分别在什么操作顺序下被触发。

澳新中高级面试考数据库，90%都落在这4个场景里。这不是新增一个独立Phase，而是Phase 3
内部、写代码之前的一个"热身"环节。

## 英语学习（并行track，2026-08-16新增）

- **现状**：TOEIC 710（2026-01），工作中英文使用比例几乎为0（0-10%）。
- **目标定位**：不设固定考试分数——新西兰AEWV签证对Software Engineer（skill level 1）
  没有任何英语测试要求（详见resumes项目`AU_NZ_Visa_Job_Notes_2026-08.md`第10节），近期
  不卡签证。聚焦实用面试表达和工作沟通能力；如果后续决定并行准备AU 482或NZ打分制路径，
  再补充正式认证考试（IELTS/PTE）。
- **战术路径（先写后说，不跳步）**：
  1. 先用英文文本提炼3个经典STAR场景（南京技术难点攻关、东京架构重构/优化、跨团队协作
     冲突），写出Situation/Task/Action/Result文本，交给AI润色成地道表达，落盘到resumes
     项目的`STAR_Stories_EN.md`，逐步建立个人词汇库。
  2. 有了打磨过的文本和词汇库打底后，再开始AI语音mock演练——不是从零裸练自由对话。
- **排期**：现在就开始，和上面的技术Phase并行推进，不是先英语后技术或先技术后英语。
- 日语：无需新增学习投入，日常工作使用即为维持N1水平的方式。

## 区域侧重

- **日本市场**：对 AWS SAA 证书的简历筛选过滤机制明显，能迅速提高通过率。
- **澳新市场**：更看重代码仓库（GitHub）的工程质量，Terraform + Docker/EKS 的真实部署代码
  非常加分；英语目前不是近期签证硬门槛（新西兰主路径无英语要求），但仍是面试/工作沟通的
  实用需求。

## 这份文件跟当前仓库的关系

- **docker-learning/**：对应 Phase 1（Docker 基础），后续可以扩展 Phase 2 的 EKS 部署练习。
- **go-learning/**：对应 Phase 3 的 Go（Gin/Echo + sqlc + PostgreSQL）BFF 旗舰项目——不再是
  "优先级低于Docker/AWS SAA的可选加分项"，而是Phase 3本身的实现载体。

---
*更新时间：2026-08-16，完整调研依据见 `~/personal/resumes/Backend_Language_Strategy_2026-08.md`、
`LinkedIn_Job_Leads_2026-08.md`、`AU_NZ_Visa_Job_Notes_2026-08.md` 与 `STAR_Stories_EN.md`。*
