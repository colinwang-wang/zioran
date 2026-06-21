# AI 辅助开发完整工作流程指南

> 本文以「知猿 (Zioran)」网课资源站的真实开发过程为案例，详细讲解如何使用 AI（Kiro CLI / Claude）从客户需求 PDF 开始，一步步完成需求分析、UI 设计、代码生成、测试的全流程。
>
> 面向：有基本编程概念但不熟悉 AI 辅助开发的初学者。

---

## 目录

1. [总览：AI 辅助开发的核心思路](#1-总览)
2. [环境准备](#2-环境准备)
3. [阶段一：从客户 PDF 到结构化需求](#3-阶段一从客户-pdf-到结构化需求)
4. [阶段二：需求分析与 PRD 编写](#4-阶段二需求分析与-prd-编写)
5. [阶段三：UI/交互设计](#5-阶段三ui交互设计)
6. [阶段四：技术方案设计](#6-阶段四技术方案设计)
7. [阶段五：代码生成（多专家协作）](#7-阶段五代码生成多专家协作)
8. [阶段六：测试与验收](#8-阶段六测试与验收)
9. [阶段七：部署上线](#9-阶段七部署上线)
10. [完整 Prompt 模板库](#10-完整-prompt-模板库)
11. [常见问题与技巧](#11-常见问题与技巧)

---

## 1. 总览

### AI 辅助开发 ≠ AI 替代开发

核心理念：**人做决策，AI 做执行**。

```
你的角色：产品经理 + 架构师 + 验收人
AI的角色：分析师 + 设计师 + 程序员 + 测试员
```

### 完整工作流一览

```
客户需求 PDF
    │
    ▼ ① AI 解读 PDF，提取需求要点
结构化需求清单
    │
    ▼ ② AI 生成 PRD 文档（按模块拆分）
PRD 文档集
    │
    ▼ ③ AI 生成 UI 原型 + 设计系统
HTML 原型 + Design Token
    │
    ▼ ④ AI 出技术方案（数据库/API/架构）
技术设计文档
    │
    ▼ ⑤ AI 多专家并行写代码
可运行的项目代码
    │
    ▼ ⑥ AI 测试 + 人工验收
测试报告 + 修复
    │
    ▼ ⑦ 部署上线
生产环境
```

---

## 2. 环境准备

### 2.1 工具清单

| 工具 | 用途 | 安装 |
|------|------|------|
| Kiro CLI | AI 对话终端（核心工具） | `npm i -g @anthropic/kiro-cli` |
| VS Code / Cursor | 代码编辑器 | 官网下载 |
| Git | 版本控制 | `brew install git` |
| Node.js 18+ | 前端运行环境 | `brew install node` |
| Go 1.21+ | 后端运行环境 | `brew install go` |
| Docker | 容器化部署 | Docker Desktop |
| pnpm | 包管理器 | `npm i -g pnpm` |

### 2.2 项目初始化

```bash
# 创建项目目录
mkdir my-project && cd my-project
git init

# 创建标准目录结构
mkdir -p docs/prd docs/prototype
mkdir -p .skills .commander/{prompts,status,contracts,phases,system-prompts}

# 把客户给的需求PDF放进来
cp ~/Downloads/客户需求.pdf docs/
```

### 2.3 关键文件说明

```
my-project/
├── docs/
│   ├── 客户需求.pdf           ← 客户原始需求文件
│   ├── prd/                   ← AI 生成的 PRD 文档
│   ├── prototype/             ← AI 生成的 HTML 原型
│   └── DESIGN.md              ← AI 生成的设计系统
├── .skills/                   ← AI 专家的行为规范（SKILL 文件）
├── .commander/                ← 多专家协作调度文件
├── CLAUDE.md                  ← 项目入口说明（AI 读取的第一份文件）
└── CONTEXT.md                 ← 项目上下文（部署信息、密码等）
```

---

## 3. 阶段一：从客户 PDF 到结构化需求

### 3.1 目标

把客户给的 PDF（通常是思维导图、Word、草图）转化为 AI 能精确理解的结构化文本。

### 3.2 实际操作步骤

**第一步：让 AI 读取 PDF 并提取内容**

打开终端，进入项目目录：

```bash
cd my-project
kiro chat
```

输入 Prompt：

```
请阅读 docs/网站开发 功能需求 思维导图.pdf，提取其中的所有功能需求，
按以下格式输出：

1. 项目基本信息（名称、域名、定位）
2. 用户角色（有哪些角色，各自能做什么）
3. 功能模块清单（按优先级排列）
4. 每个模块的核心功能点
5. 设计要求（颜色、风格、参考网站）
6. 技术偏好（如果有）
```

**第二步：补充参考网站分析（如客户给了参考）**

```
请分析参考网站 https://chenwenb.com/ ，提取以下信息：
1. 网站的完整页面结构（所有页面 URL + 功能）
2. 全局组件（导航栏、底部、搜索等）
3. 核心业务流程（用户从注册到完成购买的完整路径）
4. 付费模式（VIP、单买等）
5. 技术栈（从页面源码中分析）

输出格式：Markdown PRD 格式，保存到 docs/chenwenb-prd.md
```

**第三步：人工审查与决策**

AI 提取完后，你需要做的决策：
- 哪些功能是第一期做的？哪些可以放后面？
- 参考网站的哪些功能不要？哪些需要改？
- 有没有客户没说但必须有的功能（如隐私政策）？

```
根据以上需求分析，我做以下决策：
1. 一期只做：用户系统、课程浏览/搜索、付费（VIP+单买）、留言板
2. 二期做：投稿、问答社区、工单系统
3. 技术栈确定：Go(Gin) 后端 + Next.js 14 前端 + React Ant Design 管理后台
4. 域名：zioran.com，站名：知猿
5. 主色：#ff0036（红色系）

请基于以上决策，输出一份完整的项目概览文档，保存到 docs/prd/00-overview.md
```

### 3.3 本阶段产出

```
docs/
├── 网站开发 功能需求 思维导图.pdf    ← 客户原始文件
├── chenwenb-prd.md                  ← 参考网站分析
└── prd/
    └── 00-overview.md               ← 项目概览（你审批过的版本）
```

### 3.4 关键原则

> ⚠️ **不要跳过人工审查**。AI 提取的内容可能有遗漏或误解。你必须逐条确认，做出优先级决策后再让 AI 继续。

---

## 4. 阶段二：需求分析与 PRD 编写

### 4.1 目标

将概览文档拆分为多个模块的详细 PRD（产品需求文档），每个模块包含完整的用户故事、功能点、业务规则。

### 4.2 按模块逐个生成 PRD

**Prompt 模板：**

```
基于 docs/prd/00-overview.md 的项目概览，
请编写「用户系统」模块的详细 PRD，包含：

1. 模块概述（一句话说清模块做什么）
2. 用户角色与权限矩阵
3. 功能清单（编号，方便后续追踪）
4. 每个功能的详细描述：
   - 触发条件
   - 正常流程（步骤 1-2-3）
   - 异常流程（密码错误、验证码过期等）
   - 业务规则（密码长度限制、手机号格式等）
5. 接口预览（Method + Path + 简要说明）
6. 非功能需求（性能、安全）

保存到 docs/prd/01-user.md
```

**重复上述过程，按模块生成：**

| 文件 | 模块 | 要点 |
|------|------|------|
| `01-user.md` | 用户系统 | 注册、登录、个人中心、密码找回 |
| `02-course.md` | 课程模块 | 分类、列表、详情、搜索、标签 |
| `03-payment.md` | 付费系统 | 金币充值、VIP购买、单课购买 |
| `04-community.md` | 社区功能 | 留言板、评论、投稿 |
| `05-frontend.md` | 前端设计 | 页面路由、SEO、响应式 |
| `06-admin.md` | 管理后台 | 仪表盘、内容管理、用户管理 |
| `07-database.md` | 数据库设计 | ER图、表结构、索引 |
| `08-api.md` | API 接口 | 统一响应格式、认证、限流 |
| `09-pages.md` | 页面细节 | 每个页面的组件组成 |

### 4.3 PRD 质量检查 Prompt

每个 PRD 写完后，用这个 Prompt 做质量检查：

```
请检查 docs/prd/01-user.md 是否存在以下问题：
1. 是否有功能描述但没有异常处理说明？
2. 是否有两个功能之间有逻辑矛盾？
3. 是否有含糊的描述（如"合适的"、"足够的"）需要量化？
4. 接口设计是否覆盖了所有功能点？
5. 安全相关的功能是否有明确的防护措施说明？

列出所有问题，并给出修改建议。
```

### 4.4 本阶段产出

```
docs/prd/
├── 00-overview.md
├── 01-user.md
├── 02-course.md
├── 03-payment.md
├── 04-community.md
├── 05-frontend.md
├── 06-admin.md
├── 07-database.md
├── 08-api.md
└── 09-pages.md
```

---

## 5. 阶段三：UI/交互设计

### 5.1 目标

生成可交互的 HTML 原型页面 + 统一的设计系统（Design Token）。

### 5.2 设计系统生成

**Prompt：**

```
请基于以下信息生成设计系统文档（Design Token）：
- 参考网站：https://chenwenb.com/
- 主色：#ff0036
- 风格：干净、专业、内容为主
- 目标用户：设计师、创作者

请输出 YAML 格式的 Design Token，包含：
1. 完整色彩系统（主色、辅助色、中性色、语义色、暗色模式）
2. 字体系统（标题层级、正文、辅助文字）
3. 间距系统（4px 基准网格）
4. 圆角规范
5. 阴影层级
6. 断点定义（响应式）
7. 组件变体（按钮、卡片、输入框、标签等）

保存到 docs/DESIGN.md
```

### 5.3 HTML 原型生成

**逐页面生成原型：**

```
请基于 docs/prd/09-pages.md 中首页的描述和 docs/DESIGN.md 的设计系统，
生成一个可直接在浏览器打开的 HTML 原型页面。

要求：
1. 使用内联 CSS（单文件，无外部依赖）
2. 包含所有首页组件：Header、金刚区、Banner、搜索、最新发布、课程列表、VIP区、Footer
3. 使用真实的中文占位内容（不要 Lorem ipsum）
4. 响应式布局（桌面 + 移动端）
5. 交互状态用 CSS hover 效果展示
6. 颜色/字体严格使用 DESIGN.md 定义的值

保存到 docs/prototype/index.html
```

**生成其他页面：**

```
继续生成以下页面的 HTML 原型：
1. docs/prototype/courses.html — 课程列表页（含分类筛选）
2. docs/prototype/detail.html — 课程详情页（含购买/下载按钮）
3. docs/prototype/login.html — 登录注册页
4. docs/prototype/user.html — 个人中心
5. docs/prototype/vip.html — VIP 介绍页
6. docs/prototype/guestbook.html — 留言板
7. docs/prototype/admin.html — 管理后台

每个页面都需要：内联CSS、真实占位数据、响应式、与 index.html 风格统一
```

### 5.4 原型审查

在浏览器中逐个打开原型文件：

```bash
open docs/prototype/index.html
```

审查要点：
- 布局是否符合 PRD 描述？
- 配色是否一致？
- 移动端（F12 → 手机模式）展示是否正常？
- 交互逻辑是否清晰？

发现问题就告诉 AI 修改：

```
docs/prototype/index.html 有以下问题：
1. 金刚区图标太大，应该是 48x48px
2. 课程卡片之间间距太小，需要 16px
3. Header 在移动端没有汉堡菜单
请修改。
```

### 5.5 本阶段产出

```
docs/
├── DESIGN.md                 ← 设计系统（Design Token）
└── prototype/
    ├── index.html            ← 首页原型
    ├── courses.html          ← 课程列表
    ├── detail.html           ← 课程详情
    ├── login.html            ← 登录注册
    ├── user.html             ← 个人中心
    ├── vip.html              ← VIP页
    ├── guestbook.html        ← 留言板
    └── admin.html            ← 管理后台
```

---

## 6. 阶段四：技术方案设计

### 6.1 目标

在写代码之前，确定数据库结构、API 接口契约、项目架构。

### 6.2 数据库设计

**Prompt：**

```
基于 docs/prd/ 下所有PRD文档的业务需求，设计完整的数据库结构。

要求：
1. 列出所有表（中文名 + 英文表名）
2. 每张表的字段（名称、类型、约束、说明）
3. 表间关系（外键、关联说明）
4. 必要的索引设计
5. 种子数据说明（初始分类、管理员账号等）

输出SQL迁移文件，保存到 backend/migrations/001_init_schema.sql
同时更新 docs/prd/07-database.md 加入ER图描述
```

### 6.3 API 接口契约

**Prompt：**

```
基于 docs/prd/ 下的所有需求文档和数据库设计，
生成完整的 API 接口文档。

要求：
1. RESTful 风格
2. 统一响应格式：{"code": 0, "message": "ok", "data": {...}}
3. 分页格式：{"items": [], "total": 100, "page": 1, "pageSize": 20}
4. 统一错误码定义（40001参数错误、40101未认证、40301无权限等）
5. 每个接口包含：Method、Path、请求参数、响应格式、错误场景
6. 按模块分组（用户、课程、支付、管理员）
7. 认证方式说明（JWT Bearer Token）

保存到 docs/prd/08-api.md
```

### 6.4 项目架构设计

**Prompt：**

```
基于已确定的技术栈和API设计，输出项目目录结构和架构说明。

后端 (Go + Gin)：
- 分层架构：api(handler) → service → repository → model
- 统一中间件：JWT认证、错误处理、日志、CORS
- 配置管理：yaml文件，支持环境变量覆盖

前端 (Next.js 14)：
- App Router 目录结构
- API 封装层（SSR 和客户端分别处理）
- 组件划分（布局组件、业务组件、通用组件）

管理后台 (React + Ant Design)：
- 路由结构
- API 封装
- 页面组件划分

保存到 CONTEXT.md
```

### 6.5 本阶段产出

```
├── CONTEXT.md                      ← 项目架构与上下文
├── backend/migrations/
│   └── 001_init_schema.sql         ← 数据库建表SQL
└── docs/prd/
    ├── 07-database.md              ← 数据库设计文档
    └── 08-api.md                   ← API 接口契约
```

---

## 7. 阶段五：代码生成（多专家协作）

### 7.1 核心理念：多专家并行

这是最核心的阶段。我们使用「多专家协作」模式：开多个终端窗口，每个窗口让 AI 扮演不同的专家角色，由你（或一个「指挥官」AI）来协调。

```
┌───────────────────────────────────────────────────┐
│  终端 1：指挥官（协调分配任务）                      │
├───────────────────────────────────────────────────┤
│  终端 2：后端专家（Go API 开发）                    │
│  终端 3：前端专家（Next.js 页面开发）               │
│  终端 4：管理端专家（React 管理后台）               │
│  终端 5：测试专家（接口测试 + UI 走查）             │
└───────────────────────────────────────────────────┘
```

### 7.2 准备 SKILL 文件（专家行为规范）

每个专家需要一份 SKILL 文件，定义它的行为边界和代码规范：

**`.skills/go-backend/SKILL.md` 示例关键内容：**

```markdown
---
name: go-backend
description: Go 后端 API 开发规范
---

# Go 后端开发规范

## 代码结构
- Handler 层：参数绑定 + 调用 Service + 返回响应
- Service 层：业务逻辑 + 事务管理
- Repository 层：数据库操作（GORM）
- Model 层：数据结构定义

## 命名规范
- 文件名：snake_case
- 结构体：PascalCase
- JSON 字段：snake_case
- 常量：全大写 + 下划线

## 错误处理
- 不 panic，所有错误必须返回
- 使用统一错误码
- 敏感信息不暴露到响应中

## 交付标准
- [ ] go build 零错误
- [ ] 所有接口手动验证返回正确业务码
- [ ] 数据库迁移可重复执行
```

### 7.3 单终端模式（推荐新手）

如果觉得多终端复杂，可以用单终端让 AI 自己扮演多角色：

```bash
cd my-project
kiro chat
```

```
你是项目总指挥，遵循以下工作流：
1. 阅读 docs/prd/ 下的所有需求文档
2. 阅读 CONTEXT.md 了解项目架构
3. 按照以下顺序生成代码：
   Phase 1：后端 API（全部接口实现 + 数据库迁移）
   Phase 2：前端用户端（所有页面 + API 对接）
   Phase 3：管理后台（所有管理页面）

每个 Phase 完成后告诉我，我审查通过后再继续下一个 Phase。

从 Phase 1 开始。
```

### 7.4 多终端模式（推荐有经验者）

**终端 1 — 指挥官：**

```bash
kiro chat
```

```
你是项目指挥官，遵循 .skills/project-commander/SKILL.md。

当前阶段：Phase 1 — 后端 API 开发
需求文档：docs/prd/
数据库设计：docs/prd/07-database.md
API 契约：docs/prd/08-api.md

请为后端专家生成执行指令，写入 .commander/prompts/backend.md。
指令需包含：具体要实现的文件列表、接口清单、交付标准。
```

**终端 2 — 后端专家：**

```bash
kiro chat
```

```
你是后端开发专家，遵循 .skills/go-backend/SKILL.md。
请读取 .commander/prompts/backend.md 获取当前任务，开始执行。
完成后将状态报告写入 .commander/status/backend.md。
```

**终端 3 — 前端专家（等 Phase 1 完成后启动）：**

```bash
kiro chat
```

```
你是前端开发专家，遵循 .skills/web-frontend/SKILL.md。
请读取 .commander/prompts/frontend.md 获取当前任务，开始执行。
参考 docs/prototype/ 下的 HTML 原型实现页面。
完成后将状态报告写入 .commander/status/frontend.md。
```

### 7.5 代码生成的关键 Prompt 技巧

**技巧 1：给具体文件路径**

❌ 不好：`请实现用户注册功能`
✅ 好：

```
请实现用户注册功能，涉及以下文件：
- backend/internal/model/user.go — User 结构体
- backend/internal/repository/user_repo.go — CreateUser, FindByPhone
- backend/internal/service/user_service.go — Register (含密码加密、验证码校验)
- backend/internal/api/user_handler.go — RegisterHandler
- backend/internal/api/router.go — 注册路由 POST /api/v1/register

业务规则参考：docs/prd/01-user.md 第3节
API 格式参考：docs/prd/08-api.md
```

**技巧 2：分层实现，逐步验证**

不要让 AI 一次性写完所有代码。正确的节奏：

```
第一步：请先实现数据库迁移和 Model 定义
→ 审查通过

第二步：请实现 Repository 层（所有 CRUD 操作）
→ 审查通过

第三步：请实现 Service 层（业务逻辑）
→ 审查通过

第四步：请实现 Handler 层和路由注册
→ 审查通过

第五步：请启动服务并用 curl 测试各接口
→ 验证通过 → Phase 完成
```

**技巧 3：参考现有代码风格**

```
请参考 backend/internal/api/user_handler.go 的代码风格，
实现 course_handler.go 的所有接口。
保持相同的错误处理模式、响应格式和注释风格。
```

### 7.6 本阶段产出

```
backend/
├── cmd/server/main.go
├── internal/
│   ├── api/          ← 路由 + Handler
│   ├── service/      ← 业务逻辑
│   ├── repository/   ← 数据访问
│   ├── model/        ← 数据模型
│   └── middleware/   ← 中间件
├── pkg/              ← 公共包（短信、支付、配置）
├── migrations/       ← SQL 迁移
└── config.yaml

frontend/
├── src/app/          ← 页面路由（Next.js App Router）
├── src/components/   ← 通用组件
├── src/lib/          ← API 封装
├── src/types/        ← TypeScript 类型
└── src/contexts/     ← 状态管理

admin/
├── src/pages/        ← 管理页面
├── src/api/          ← API 封装
├── src/layouts/      ← 布局组件
└── src/types/        ← 类型定义
```

---

## 8. 阶段六：测试与验收

### 8.1 API 接口测试

**Prompt：**

```
请对以下后端 API 进行完整测试：

1. 启动后端服务
2. 逐个接口发送请求（用 curl 或 HTTP 测试工具）
3. 验证：
   - HTTP 状态码正确
   - 业务响应码正确（code: 0 表示成功）
   - 返回数据结构与 docs/prd/08-api.md 一致
   - 异常场景返回正确错误码
4. 按以下格式输出测试报告：

| 接口 | Method | 状态 | 备注 |
|------|--------|------|------|
| /api/v1/register | POST | ✅ PASS | - |
| /api/v1/login | POST | ❌ FAIL | 错误码返回格式不对 |

保存报告到 .commander/status/qa.md
```

### 8.2 前后端联调测试

```
请执行前后端联调测试：

1. 同时启动后端（端口 8080）和前端（端口 3000）
2. 在浏览器中完成以下用户路径：
   - 注册新用户 → 登录 → 浏览课程 → 查看详情 → 充值金币 → 购买课程
3. 检查：
   - 前端请求路径与后端路由完全匹配
   - JSON 字段名前后端一致（不能前端 camelCase 后端 snake_case 不匹配）
   - 分页数据格式正确
   - Token 认证流程正常
4. 记录所有不通过的问题，输出修复指令
```

### 8.3 UI 走查

```
请对比 docs/prototype/ 下的 HTML 原型，检查前端实现：

1. 打开 docs/prototype/index.html 和实际运行的首页
2. 逐个区域对比：
   - Header 布局和导航项一致
   - 金刚区图标数量和排列一致
   - 课程卡片样式和信息展示一致
   - 颜色是否使用了 DESIGN.md 定义的值
3. 检查响应式（桌面 / 平板 / 手机三个断点）
4. 输出差异清单和修改建议
```

### 8.4 Bug 修复循环

发现 Bug 后的标准流程：

```
测试发现以下问题：

1. [BUG-001] 登录接口返回的 token 字段名是 access_token，
   但前端代码读取的是 token，导致登录后立即掉线
   → 修复方案：统一为 token

2. [BUG-002] 课程列表接口返回 {list: [...]}，
   但前端期望 {items: [...]}
   → 修复方案：后端改为 items（与 API 文档一致）

请逐个修复以上问题，修复后重新运行受影响的测试用例验证。
```

### 8.5 本阶段产出

```
.commander/status/
├── qa.md                    ← 测试报告
├── backend.md               ← 后端完成状态
├── frontend.md              ← 前端完成状态
└── admin.md                 ← 管理端完成状态
```

---

## 9. 阶段七：部署上线

### 9.1 Docker 化

**Prompt：**

```
请为项目生成 Docker 部署配置：

1. backend/Dockerfile — Go 编译 + 最小运行镜像
2. frontend/Dockerfile — Next.js 构建 + standalone 运行
3. admin/Dockerfile — Vite 构建 + Nginx 静态托管
4. docker-compose.yml — 编排所有服务（含 MySQL、Redis、Nginx）
5. nginx/nginx.conf — 反向代理配置：
   - / → frontend:3000
   - /api/ → backend:8080
   - /admin/ → admin 静态文件

要求：
- 生产环境优化（多阶段构建、最小镜像）
- 环境变量通过 .env 文件注入
- 数据持久化（MySQL 数据卷）
```

### 9.2 一键部署脚本

```
请生成 Makefile 包含以下命令：

make deploy-full     # SSH到服务器，拉代码，重建所有容器
make quick-deploy    # 只重建后端
make status          # 查看服务状态
make logs-api        # 查看后端日志
make db-backup       # 备份数据库
make health          # 健康检查（curl 各服务）
```

### 9.3 上线前检查清单

```
请生成上线前的 checklist，检查以下项目：
1. 测试万能码（验证码 0000、短信 000000）是否已移除或改为环境变量控制
2. 数据库密码是否已从代码中移除，改为环境变量
3. CORS 配置是否只允许正式域名
4. 日志级别是否切换到 production
5. 所有 TODO/MOCK 占位是否已清理
6. SSL 证书是否配置
7. 数据库备份策略是否就绪
```

---

## 10. 完整 Prompt 模板库

### 10.1 需求分析阶段

```
【PDF 解读模板】
请阅读 {文件路径}，提取所有功能需求，按以下格式输出：
1. 项目基本信息
2. 用户角色
3. 功能模块清单
4. 设计要求
5. 技术偏好
不要遗漏任何细节，如有不确定的地方请标注 [待确认]。
```

```
【参考网站分析模板】
请分析参考网站 {URL}，输出完整的功能 PRD：
1. 页面结构（URL + 功能）
2. 全局组件
3. 核心业务流程
4. 付费模式
5. 技术栈分析
保存到 docs/{filename}.md
```

### 10.2 PRD 编写阶段

```
【模块 PRD 模板】
基于 docs/prd/00-overview.md，编写「{模块名}」的详细 PRD：
1. 模块概述
2. 用户角色与权限
3. 功能清单（编号）
4. 每个功能的详细描述（正常流程 + 异常流程 + 业务规则）
5. 接口预览
6. 非功能需求
保存到 docs/prd/{编号}-{模块名}.md
```

### 10.3 设计阶段

```
【设计系统模板】
请基于主色 {色值}、风格 {描述}、参考 {网站}，
生成完整设计系统（YAML 格式）：
colors / typography / spacing / radius / shadow / breakpoints / components
保存到 docs/DESIGN.md
```

```
【HTML 原型模板】
请基于 docs/prd/09-pages.md 中 {页面名} 的描述，
结合 docs/DESIGN.md 的设计系统，
生成单文件 HTML 原型（内联CSS、真实中文数据、响应式）。
保存到 docs/prototype/{filename}.html
```

### 10.4 代码生成阶段

```
【后端开发模板】
请实现 {模块名} 的后端 API，涉及文件：
- model: backend/internal/model/{name}.go
- repository: backend/internal/repository/{name}_repo.go
- service: backend/internal/service/{name}_service.go
- handler: backend/internal/api/{name}_handler.go
- router: 注册到 backend/internal/api/router.go

业务规则参考：docs/prd/{编号}.md
API 格式参考：docs/prd/08-api.md
代码风格参考：{已有文件}
```

```
【前端开发模板】
请实现 {页面名} 页面：
- 路由文件：frontend/src/app/{path}/page.tsx
- 布局参考：docs/prototype/{name}.html
- API 数据：使用 frontend/src/lib/services.ts 中的方法
- 组件复用：参考 frontend/src/components/ 已有组件
- 样式：Tailwind CSS，颜色使用 DESIGN.md 定义的变量
```

### 10.5 测试阶段

```
【接口测试模板】
请对 {模块} 相关的所有 API 执行测试：
1. 正常路径（正确参数 → 期望响应）
2. 参数校验（缺少必填字段 → 40001 错误）
3. 权限控制（未登录 → 40101，无权限 → 40301）
4. 业务异常（重复注册 → 对应错误码）

输出格式：表格形式的测试报告
```

---

## 11. 常见问题与技巧

### Q1：AI 一次性输出太多代码，质量不高怎么办？

**答：分步走。** 不要一次让 AI 写完整个项目。按照 Model → Repository → Service → Handler 的顺序逐层实现，每层写完后审查通过再继续。

### Q2：前后端接口对不上怎么办？

**答：契约优先（Contract First）。** 先定义好 API 文档，后端和前端都严格按文档实现。联调时发现不匹配，以文档为准修改代码。

### Q3：AI 生成的代码风格不统一怎么办？

**答：用 SKILL 文件约束 + 给参考代码。** 让 AI「参考 xxx 文件的风格」比抽象描述有效得多。

### Q4：AI 做了不该做的事（多加功能、改架构）怎么办？

**答：明确边界。** 在 Prompt 中写清楚：
- 「只修改以下文件：...」
- 「不要添加任何文档中没有描述的功能」
- 「保持现有架构不变」

### Q5：多个 AI 专家之间代码冲突怎么办？

**答：文件隔离 + 契约通信。**
- 后端专家只改 `backend/` 目录
- 前端专家只改 `frontend/` 目录
- 通信通过 API 契约文件 + `.commander/` 下的指令/状态文件

### Q6：代码能编译但功能有 Bug 怎么办？

**答：要求 AI 做端到端验证。** 不要只问「实现完了吗」，要问「请启动服务，用 curl 调用接口验证返回是否正确」。

### Q7：如何节省 AI 调用成本？

**答：混合模型策略。**
- 需求分析、架构设计用强模型（Claude Opus/Sonnet）
- 写具体代码用普通模型（GPT-4o / DeepSeek）
- 指挥官生成的指令越详细，弱模型执行的代码质量越高

### Q8：项目做到一半需求变了怎么办？

**答：从文档开始改。**
1. 先更新 PRD 文档
2. 再更新 API 契约
3. 然后让 AI 根据新文档修改代码
4. 最后重新测试

**永远不要直接改代码而不更新文档。**

---

## 附录 A：本项目（知猿）的完整文件清单

```
zioran/
├── CLAUDE.md                          # 项目入口（AI 首先阅读）
├── CONTEXT.md                         # 完整项目上下文（技术栈、部署信息）
├── WORKFLOW.md                        # 多专家协作工作流说明
├── Makefile                           # 运维命令集
├── docker-compose.yml                 # 容器编排
│
├── docs/
│   ├── 网站开发 功能需求 思维导图.pdf  # 客户原始需求
│   ├── chenwenb-prd.md               # 参考网站分析
│   ├── DESIGN.md                     # 设计系统（700+行）
│   ├── project_context.md            # 项目简介
│   ├── fix-list.md                   # 问题修复清单
│   ├── prd/                          # PRD 文档集（10份）
│   │   ├── 00-overview.md
│   │   ├── 01-user.md
│   │   ├── 02-course.md
│   │   ├── 03-payment.md
│   │   ├── 04-community.md
│   │   ├── 05-frontend.md
│   │   ├── 06-admin.md
│   │   ├── 07-database.md
│   │   ├── 08-api.md
│   │   └── 09-pages.md
│   └── prototype/                    # HTML 原型（8页）
│       ├── index.html
│       ├── courses.html
│       ├── detail.html
│       ├── login.html
│       ├── user.html
│       ├── vip.html
│       ├── guestbook.html
│       └── admin.html
│
├── .skills/                          # AI 专家行为规范
│   ├── project-commander/SKILL.md    # 项目指挥官
│   ├── fullstack-planning/SKILL.md   # 全栈规划
│   ├── go-backend/SKILL.md           # Go 后端规范
│   ├── web-frontend/SKILL.md         # 前端规范
│   ├── web-admin/SKILL.md            # 管理端规范
│   ├── qa-testing/SKILL.md           # 测试规范
│   ├── coding-guidelines/SKILL.md    # 编码通用规范
│   └── test-driven-development/SKILL.md  # TDD 规范
│
├── .commander/                       # 多专家协作调度
│   ├── prompts/                      # 指挥官→专家 的指令
│   ├── status/                       # 专家→指挥官 的状态报告
│   ├── contracts/                    # API 契约文件
│   ├── phases/                       # 阶段记录
│   └── system-prompts/               # 各角色系统提示词
│
├── backend/                          # Go 后端（77个API路由）
├── frontend/                         # Next.js 用户端（19个页面）
├── admin/                            # React 管理后台（14个页面）
└── nginx/                            # Nginx 配置
```

---

## 附录 B：从零到上线的时间线参考

| 阶段 | 耗时 | 人工工作 | AI 工作 |
|------|------|----------|---------|
| 需求提取 | 2小时 | 审查确认、优先级决策 | 读PDF、分析参考网站 |
| PRD 编写 | 3小时 | 审查每份PRD、补充细节 | 生成10份PRD文档 |
| UI 设计 | 2小时 | 浏览器审查原型、提修改 | 生成设计系统+8页原型 |
| 技术方案 | 1小时 | 确认架构选型 | 数据库设计、API契约 |
| 后端开发 | 4小时 | 审查代码、验证接口 | 生成全部后端代码 |
| 前端开发 | 4小时 | 审查页面、提UI修改 | 生成全部前端页面 |
| 管理后台 | 2小时 | 审查功能完整性 | 生成管理后台代码 |
| 测试修复 | 3小时 | 验收测试报告 | 执行测试+修复Bug |
| 部署上线 | 1小时 | 确认环境配置 | 生成Docker+部署脚本 |
| **合计** | **~22小时** | **~8小时**（决策+审查） | **~14小时**（执行） |

> 💡 传统方式预估：一个全栈开发者约需 3-4 周。AI 辅助将实际开发时间压缩到 2-3 天。

---

## 附录 C：推荐阅读

- 项目工作流详解：`WORKFLOW.md`
- 多专家系统提示词：`.commander/system-prompts/`
- 专家能力规范：`.skills/*/SKILL.md`
- 实际测试报告案例：`.commander/status/qa.md`
