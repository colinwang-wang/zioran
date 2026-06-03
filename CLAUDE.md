# 知猿 (zioran)

网课资源付费下载站。Go(Gin) + Next.js 14 前后端分离架构。

## 技术栈
- 后端: Go + Gin + GORM + PostgreSQL + Redis
- 用户端: Next.js 14 (App Router) + Tailwind CSS + shadcn/ui
- 管理后台: React + TypeScript + Ant Design 5
- 支付: 微信支付 + 支付宝
- 部署: Docker Compose

## 项目结构
```
zioran/
├── docs/prd/          # 产品需求文档
├── .skills/           # 契约开发规范
├── .commander/        # 多专家协作协议
├── WORKFLOW.md        # 协作工作流说明
├── backend/           # Go 后端 (待创建)
├── frontend/          # Next.js 用户端 (待创建)
└── admin/             # React 管理后台 (待创建)
```

## Skills
- `.skills/project-commander/` — 项目总指挥规范
- `.skills/fullstack-planning/` — 全栈项目规划
- `.skills/go-backend/` — Go 后端 API 开发规范
- `.skills/web-frontend/` — Web 用户端规范
- `.skills/web-admin/` — Web 管理端规范
- `.skills/qa-testing/` — 测试专家规范
- `.skills/coding-guidelines/` — 编码通用规范

## 快速开始
```bash
# 指挥官模式（推荐）
kiro chat
# 输入: 你是项目总指挥，遵循 .skills/project-commander/SKILL.md 规范。请分析 docs/prd/ 下的需求文档，开始 Phase 1。
```
