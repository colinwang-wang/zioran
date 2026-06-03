# 多专家协作工作流

## 架构概览

```
zioran/
├── .commander/
│   ├── prompts/              # 指挥官生成的指令文件（专家读取）
│   │   ├── backend.md        # 当前给后端专家的指令
│   │   ├── admin.md          # 当前给管理端专家的指令
│   │   └── miniapp.md        # 当前给小程序专家的指令
│   ├── contracts/            # 接口契约（唯一真相来源）
│   │   └── api-v1.yaml       # OpenAPI 契约文件
│   ├── status/               # 各专家状态报告
│   │   ├── backend.md        # 后端专家完成报告
│   │   ├── admin.md          # 管理端专家完成报告
│   │   └── miniapp.md        # 小程序专家完成报告
│   ├── phases/               # 阶段记录
│   │   └── phase-01.md       # 当前阶段目标与进度
│   └── system-prompts/       # 各角色系统提示词
│       ├── commander.md
│       ├── backend.md
│       ├── admin.md
│       ├── miniapp.md
│       └── qa.md
├── .skills/                  # 各角色 SKILL 规范
└── CLAUDE.md                 # 项目总览（可选）
```

## 角色分工

| 终端 | 角色 | 职责 |
|------|------|------|
| 终端 1 | 指挥官 (Commander) | 需求拆解、生成指令、验收产出、分派 bug |
| 终端 2 | 后端专家 (Backend) | Go + Gin + PostgreSQL API 开发 |
| 终端 3 | 前端专家 (Frontend) | Next.js 14 用户端网站开发 |
| 终端 4 | 管理端专家 (Admin) | React + Ant Design 管理后台 |
| 终端 5 | 测试专家 (QA) | 接口测试、UI 走查、联调验证、出测试报告 |

## 工作流循环

```
┌─────────────────────────────────────────────────┐
│  指挥官                                          │
│  1. 分析需求 / 验收上轮产出                       │
│  2. 写入 .commander/prompts/{role}.md            │
│  3. 通知专家："指令已更新，开始执行"               │
└──────────────────────┬──────────────────────────┘
                       │
         ┌─────────────┼─────────────┐
         ▼             ▼             ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│ 后端专家      │ │ 前端专家      │ │ 管理端专家    │
│ 读取 prompt  │ │ 读取 prompt  │ │ 读取 prompt  │
│ 执行任务     │ │ 执行任务     │ │ 执行任务     │
│ 写入 status  │ │ 写入 status  │ │ 写入 status  │
└──────────────┘ └──────────────┘ └──────────────┘
         │             │             │
         └─────────────┼─────────────┘
                       ▼
              指挥官验收 → 调度测试
                       │
                       ▼
              ┌──────────────┐
              │ 测试专家      │
              │ 执行测试     │
              │ 写入测试报告  │
              └──────┬───────┘
                     ▼
              指挥官读取测试报告
              ├── 全部 PASS → 下一个 Phase
              └── 有 BUG → 生成修复指令 → 专家修复 → 回归验证
```

## 指令文件格式

### prompts/{role}.md

```markdown
# Phase {N} — {角色} 指令

> 状态: PENDING | IN_PROGRESS | DONE
> 依赖: 无 | 等待 {role} 完成
> 更新时间: {timestamp}

## 背景
为什么要做这件事

## 任务列表
1. 具体任务（附文件路径或接口定义）
2. ...

## 交付标准
- [ ] 可验证的检查项

## 参考
- 契约文件: .commander/contracts/api-v1.yaml
- SKILL 规范: .skills/{role}/SKILL.md
```

### status/{role}.md

```markdown
# {角色} 状态报告

> 状态: DONE
> 完成时间: {timestamp}

## 完成内容
- 已完成的任务列表

## 自检报告
- [x] 检查项 1
- [x] 检查项 2

## 问题与阻塞
- 无 / 描述遇到的问题
```

## 操作指南

### 方式一：单终端 subagent 模式（推荐）

```bash
cd zioran
kiro chat
```

```
你是项目总指挥，遵循 .skills/project-commander/SKILL.md 规范。
产品文档在 docs/，技术栈 Go+Gin / React+AntD / 微信小程序。
请自动推进：需求分析 → 契约 → 开发 → 测试 → 部署。
```

### 方式二：多终端手动模式

```bash
# 终端 1 — 指挥官
kiro chat --system-prompt .commander/system-prompts/commander.md

# 终端 2 — 后端专家
kiro chat --system-prompt .commander/system-prompts/backend.md

# 终端 3 — 前端专家
kiro chat --system-prompt .commander/system-prompts/frontend.md

# 终端 4 — 管理端专家
kiro chat --system-prompt .commander/system-prompts/admin.md

# 终端 5 — 测试专家
kiro chat --system-prompt .commander/system-prompts/qa.md
```

专家启动后告诉它：
```
读取 .commander/prompts/{role}.md 开始执行
```

### 方式三：混合模型模式（省成本）

- 指挥官用强模型（Claude Opus / Sonnet）
- 专家用便宜模型（DeepSeek / 本地模型）
- 指挥官生成的指令足够详细，弱模型也能按指令完成

## 文件交互协议

```
指挥官写指令 → .commander/prompts/{role}.md
专家读指令   ← .commander/prompts/{role}.md
专家写报告   → .commander/status/{role}.md
指挥官读报告 ← .commander/status/{role}.md
契约文件     ↔ .commander/contracts/api-v1.yaml
```

## 契约管理
- 契约文件放在 `.commander/contracts/`
- 后端专家负责生成/更新契约
- 前端专家只读契约，不修改
- 指挥官审核契约变更
