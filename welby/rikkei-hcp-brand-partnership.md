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

## 公司背景（官网核实，2026-09-09 查证 https://rikkeisoft.com）

官网"会社概要"页面显示，Rikkei 实际是**母公司（越南）+ 日本法人**两层
结构：

### 越南母公司：株式会社Rikkeisoft

| 项目 | 内容 |
|---|---|
| 设立 | 2012年4月6日（平成24年） |
| 资本金 | 6亿日元 |
| 代表 | 代表取締役会长 Ta Son Tung |
| 员工数 | 2,368名（2024年12月现在） |
| 总部 | **河内**（Vinacomin Tower, Cau Giay Dist.） |
| 分社 | 岘港（ダナン）、顺化（フエ）、胡志明市 |

主营业务起家于软件外包/离岸开发（オフショア開発），后续扩展方向包括
另有独立法人 **株式会社Rikkei AI**（河内），说明集团也在往 AI 方向布局。

### 日本法人：株式会社リッケイ

| 项目 | 内容 |
|---|---|
| 设立 | 2016年3月1日（平成28年） |
| 资本金 | 8,386万8,386日元 |
| 代表 | 代表取締役社长 Bui Quang Huy |
| 员工数 | 325名（2024年12月现在） |
| 加盟团体 | 一般社団法人 情報サービス産業協会（JISA） |
| 事业内容 | オフショア開発事业、採用支援事业、コンサルティング事业、労働者派遣事业（派 13-315973） |
| 所在地 | 东京总部（港区芝浦 msb Tamachi）、大阪/名古屋/福冈/札幌/北陆均有支社据点 |

### 全球版图（2026-09-09 查证 rikkeisoft.com/greetings/ 官网世界地图）

官网"About Us"页面的全球据点地图，用红色（Offices and delivery centers，
实体办公/交付中心）和蓝色（Clients，纯客户所在地，无办公室）区分了各地区
性质：

| 地区 | 性质 | 具体据点/说明 |
|---|---|---|
| 越南 | 🔴 办公/交付中心 | 河内（Hanoi, **HQ**）、胡志明市、岘港（Da Nang）、顺化（Hue） |
| 日本 | 🔴 办公/交付中心 | 东京、大阪、名古屋、福冈、札幌、北陆 |
| 韩国 | 🔴 办公/交付中心 | 首尔——子公司 **Ksoft**（A Rikkei Company） |
| 泰国 | 🔴 办公/交付中心 | 曼谷——子公司 **RKThailand**（A Rikkei Company） |
| 美国 | 🔴 办公/交付中心 | Plano, Texas——子公司 **RKTech**（A Rikkei Company） |
| 新加坡 | 🔵 仅客户 | 无实体办公室 |
| 澳大利亚 | 🔵 仅客户 | 无实体办公室 |
| 新西兰 | 🔵 仅客户 | 无实体办公室 |

即：**Rikkei 是横跨越南本土＋日本＋韩国＋泰国＋美国的离岸开发集团**，在
韩国/泰国/美国分别用 Ksoft/RKThailand/RKTech 这几个子品牌落地设点；新加坡/
澳洲/新西兰只是市场（有客户），并未设办公据点——地图上"红/蓝"两色的区分
标准就是"有没有实体交付团队"，不是"有没有业务往来"。

跟 Welby 谈 HealthcarePlus（HCP）品牌合作、对应 `healthcareplus-dev.rikkei.org`
这个域名的，大概率就是这个**日本法人（株式会社リッケイ）**，而不是越南
母公司直接对接——离岸开发起家的公司里，通常是日本法人负责跟日本客户/
合作方谈业务，越南母公司负责实际研发交付，这跟 Rikkei 主营的"オフショア
開発"商业模式是一致的。

**仍未核实的部分**：Welby 和 Rikkei 之间具体的合作性质（技术授权？合资？
纯外包再转售？）——官网只能确认 Rikkei 这家公司的规模和资质，无法证明
双方合同的具体条款，如果要写进正式对外材料，仍建议找业务侧同事确认最新
的合作性质。
