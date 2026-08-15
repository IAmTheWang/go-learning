# Tech Learning Roadmap（AWS / Docker / Go / PostgreSQL）

> 这份文件是从 `~/personal/resumes` 项目里沉淀的求职策略文档中抽取出来的技术学习方向指引，
> 同步放在 `go-learning/` 和 `docker-learning/` 这两个技术练习仓库里，作为"为什么学这个、
> 优先级是什么"的背景说明。完整的市场调研 / 薪资对比 / 理由细节见 resumes 项目的
> `Backend_Language_Strategy_2026-08.md` 与 `LinkedIn_Job_Leads_2026-08.md`。

## 背景

前端工程师（Vue 3 / React + TypeScript），目标市场：日本（主战场）+ 中国 + 澳洲 + 新西兰。
当前 AWS 和后端都是薄弱项，策略是优先补强 AWS/云原生，不主动深入学习传统后端语言
（尤其是 Java）。

## 综合优先级：Docker > AWS SAA > (Terraform深化+EKS) > BFF > Go > Python > Java

| 排名 | 项目 | 定位 |
|---|---|---|
| 1 | **Docker** | 最高优先级，地基 —— `docker-learning/` 对应这一步，不依赖任何 AWS 知识，可以最先独立学 |
| 2 | **AWS SAA 认证** | 第二优先级，建立 HR 筛选门槛，同时是 EKS 的前置知识（VPC/Subnet/IAM OIDC） |
| 3 | **Go** | 锦上添花，不阻塞主线 —— `go-learning/` 对应这一步，优先级低于 Docker 和 AWS SAA |
| 4 | Python | 现状够用（Flask AI 辅助维护经验打底），不额外投入 |
| 5 | Java | 不主动学，工作中遇到靠 AI 工具读懂即可 |

## 详细学习顺序（4 个 Phase，2026-08-15 优化版）

> 相比最早的"六步计划"，这一版做了三处调整：**解耦 SAA 与 EKS**（EKS 依赖的 VPC/IAM 知识
> 应该先靠 SAA 打底，而不是跟 Docker 挤在同一个月学）、**砍掉 AWS CDK**（已有 Terraform
> 真实经验，跨云通用性/认可度更高，学第二个同类 IaC 工具边际收益低）、**把 BFF 具体化成
> 一个可展示的旗舰项目**。

| Phase | 核心内容 | 交付产物 | 周期 |
|---|---|---|---|
| **1. 容器与云基础** | Docker/docker-compose 基础 + AWS SAA 刷题体系 | 考出 **AWS SAA 认证** | 1–1.5 个月 |
| **2. IaC 与云原生** | 深化 Terraform + EKS 部署（IRSA/Ingress/Deployment/Service），此时已有 SAA 打下的 VPC/IAM/OIDC 基础，理解 EKS 的网络和权限模型会顺畅很多 | Terraform 自动化拉起 AWS 资源 + 一个 EKS 部署项目 | 1 个月 |
| **3. BFF 实战闭环** | **Hono**（首选，NestJS 备选）+ Drizzle + PostgreSQL，配合 Docker 做全链路自动化交付 | 完整全栈 Demo，**带 CI/CD 自动部署到 EKS**——这是跟纯 Dev/纯 DevOps 区分开的杀手锏 | 1 个月 |
| **4.（可选）AWS DevOps 认证** | 前 3 个 Phase 顺利完成后再考虑，不强制排期 | AWS Certified DevOps Engineer 认证，针对 17-25K 档位那类明确要 CI/CD+运维深度的 JD | 视情况 |

**AWS CDK**：不列入学习目标。处理方式和 Java 一致——不主动投入学习时间，真遇到明确要 CDK
的 JD，靠已有的 IaC 概念（Terraform）+ AI 辅助现学现用即可。

**Go 的定位**：和 Docker/Kubernetes/Terraform 同一条技能树延伸（这些工具本身都是 Go 写的），
学 Go 有助于看懂工具内部逻辑和报错信息，但属于"选修加分层"，不在上面 4 个 Phase 的主线里，
**不应抢占 Docker/AWS SAA 的学习时间**。

**Java**：不列入学习计划，工作中遇到靠 AI 工具（Claude/Cursor）辅助读懂即可。

## 区域侧重

- **日本市场**：对 AWS SAA 证书的简历筛选过滤机制明显，能迅速提高通过率。
- **澳新市场**：更看重代码仓库（GitHub）的工程质量，Terraform + Docker/EKS 的真实部署代码
  非常加分。

## 这份文件跟当前仓库的关系

- **docker-learning/**：对应 Phase 1（Docker 基础），后续可以扩展 Phase 2 的 EKS 部署练习。
- **go-learning/**：对应"Go"这个可选加分项，优先级低于 Docker 和 AWS SAA —— 如果时间紧张，
  Docker/AWS SAA 的学习应该优先于继续深入这个仓库。

---
*更新时间：2026-08-15，完整调研依据见 `~/personal/resumes/Backend_Language_Strategy_2026-08.md`
与 `LinkedIn_Job_Leads_2026-08.md`。*
