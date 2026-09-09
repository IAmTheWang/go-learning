# Rikkei 与 HealthcarePlus（HCP）品牌合作关系

> 业务知识笔记 — 源自 2026-09-09 阅读 `_deploy/designDocs/20260525T214848JST-hcp-phase-1.md`
> 设计文档 + `dev-rikkei` 系列部署 commit 时的调查
> 姊妹篇：[dda-architecture-overview.md](./dda-architecture-overview.md) 讲设备厂商连携（跟
> Rikkei 这种"品牌合作方"是完全不同性质的另一类外部关系）

## 一句话结论

**Rikkei 不是设备厂商，是运营着同一套代码、另一个品牌的业务合作方。**
`welby-phr-web` 这一套 Vue 代码，通过构建时的品牌开关，同时产出两个对外
产品：

- **マイカルテ（mykarte）**——Welby 自己的品牌，`karte-develop.welby.jp`
- **ヘルスケアパスポートplus（HealthcarePlus / HCP）**——Rikkei 运营的品牌，
  目前托管在 `healthcareplus-dev.rikkei.org`（Rikkei 自己的基础设施上，
  不是 Welby 的 AWS）

## 证据：设计文档原话

`_deploy/designDocs/20260525T214848JST-hcp-phase-1.md` 第 14 行写得很直白：

> "Welby's `welby-phr-web` is a Vue 3 SPA serving the **mykarte**（マイカルテ）
> brand at `karte-develop.welby.jp/admin2`... **Rikkei serves a second brand
> of the same product, HealthcarePlus (HCP, "ヘルスケアパスポートplus")**,
> at `healthcareplus-dev.rikkei.org`. The two brands share the same Vue
> codebase — the difference is build-time mode flags
> (`vite/build/brand.ts:isHcp = mode.includes("hcp")`) that switch favicon,
> title, colors, and the LEAFLET API `appName` parameter."

## 这份设计文档在做什么：Welby 想把 HCP 收回自己运营

文档第 16 行说明了目的——Welby 想在自己的 AWS 上也搭一套
`healthcareplus-{dev,stg,prd}.welby.jp`：

> "...a stable target for QA, integration testing, and (eventually)
> production traffic **if Welby chooses to operate HCP directly rather
> than via Rikkei**."

也就是说：**现在 HCP 这个品牌的生产环境运营依赖 Rikkei 那边的基础设施，
Welby 正在评估/推进"自己也能独立运营 HCP"这件事**，而不是要跟 Rikkei
分道扬镳——第 39-42 行的 NG（Not Goal）清单里，NG6 明确写了这次 Phase 1
还是要做"Rikkei parity"（跟 Rikkei 那边行为一致）的人工冒烟测试，说明
双方产品行为目前仍需保持对齐，是并行运营，不是替换关系。

## 代码里能看到的三类实际证据

**1. 品牌切换开关**（`vite/build/brand.ts`）——同一份 Vue 代码根据构建
mode 里有没有 `hcp` 字样，切换 favicon/标题/配色/LEAFLET API 的
`appName` 参数。

**2. package.json 里专门的 Rikkei 联调构建命令**：

```json
"build:dev-rikkei:mykarte": "vue-tsc && vite build --mode mykarte-dev-rikkei",
"build:dev-rikkei:hcp": "vue-tsc && vite build --mode hcp-dev-rikkei",
```

**3. `vite.config.ts` 里专门放行 Rikkei 的域名**：

```ts
server: {
    allowedHosts: ["healthcareplus-dev.rikkei.org", "karte-dev.rikkei.org"],
},
```

**4. 一批 `dev-rikkei` 专用 AWS 环境**（commit 历史里能看到，如
`844e9cb5`/`9257f984`/`d20c73a3`）——`karte-dev-rikkei.welby.jp`、
`healthcareplus-dev-rikkei.welby.jp`，配套 GitHub Actions OIDC 部署角色
（`98b6b3f3`/`8aab48b3`），看起来是 Welby 给 Rikkei 团队专门开的一个
**开发联调环境**，让他们能直接对接 Welby 这边搭的基础设施做测试，而不是
自己另起一套。

## 什么是 OIDC——为什么 GitHub Actions 部署要用它

上面提到的 "GitHub Actions OIDC deploy role" 涉及一个具体技术概念，
顺带记一下：

**OIDC = OpenID Connect**，建立在 **OAuth 2.0** 协议之上的一层"身份验证"
标准。区别在于目的不同：

- **OAuth 2.0 本身解决的是"授权"（Authorization）**——"允许这个 App
  访问我的哪些数据/资源"（DDA 那篇笔记里讲的 Fitbit/A&D 授权就是这个）。
- **OIDC 在 OAuth 之上加了一层"身份认证"（Authentication）**——多返回
  一个 **ID Token**（一段 JWT，里面编码了"你是谁"这类身份声明，比如
  `sub`、`email`），专门用来回答"这次请求到底是谁发起的"这个问题，而
  不只是"这次请求有没有权限"。

**在这里的具体用法（GitHub Actions → AWS）**：这是 OIDC 一个非常常见
的工业场景——CI/CD 流水线要往云上部署，传统做法是把 AWS Access Key 长期
存在 GitHub Secrets 里（一旦泄露、长期有效，风险很大）。用 OIDC 之后：

```
1. GitHub Actions workflow 运行时，GitHub 自己扮演"OIDC 身份提供方"
   （OIDC Identity Provider），签发一个短期有效的 JWT token，
   里面写明"这是 Welby-Inc/welby-phr-web 仓库的某个 workflow 在跑"

2. AWS IAM 这边预先配置了信任 GitHub 这个 OIDC Identity Provider
   （对应上文提到的 "GitHub Actions OIDC deploy role"）

3. Workflow 拿着这个短期 JWT 去 AWS 换取一个临时的 IAM Role 权限
   （sts:AssumeRoleWithWebIdentity），完事这个临时权限就过期作废
```

好处：**AWS 那边完全不用存一份长期有效的密钥**，每次部署都是"当场证明
身份 → 换一个用完即焚的临时权限"，泄露风险和管理成本都比长期 Access Key
小得多。这跟 DDA 笔记里讲的"OAuth token 会过期"是同一种"用短期换长期
安全性"的设计哲学，只是 OIDC 这里换的是云资源的操作权限，不是某个病人的
健康数据读取权限。

## 公司背景（一般常识，非本仓库验证内容）

Rikkei（常见全称 Rikkei Group，其软件外包子品牌叫 RikkeiSoft）总部在
**越南河内**，主营业务起家于软件外包/离岸开发，后来业务扩展到教育培训、
数字医疗等方向。HealthcarePlus 这个品牌大概率是 Rikkei 把 Welby 的
マイカルテ技术白标（white-label）改造成自己的产品对外运营。**这部分公司
规模、具体股权/合作性质（技术授权？合资？纯外包再转售？）没有在代码或
设计文档里找到直接证据**，如果要写进正式对外材料，建议找业务侧同事确认
最新的合作性质。
