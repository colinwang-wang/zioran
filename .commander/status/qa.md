# QA 代码审查报告 — 全模块

> 审查时间: 2026-06-04T00:26
> 结论: **FAIL（7个问题需修复）**

---

## 问题清单

| # | 模块 | 文件 | 严重程度 | 问题描述 | 建议修复 |
|---|------|------|----------|----------|----------|
| 1 | 后端 | api/router.go | 🔴 高 | PRD定义 `GET /api/v1/admin/users/:id`（用户详情）未注册路由 | 添加路由 + handler |
| 2 | 后端 | api/router.go | 🟡 中 | PRD定义 `PUT /api/v1/admin/categories/:id/status`（分类上下架）未注册路由 | 添加路由 + handler |
| 3 | 后端 | api/router.go | 🟡 中 | PRD定义 `POST /api/v1/admin/courses/batch`（批量操作）未注册路由 | 添加路由 + handler |
| 4 | 后端 | service/payment.go:67 | 🟡 中 | 充值接口 MOCK 直接模拟支付成功，生产环境需要明确标注并返回支付URL | 保持MOCK但返回 `{pay_url:"mock://success"}` 以便前端逻辑完整 |
| 5 | 后端 | service/auth.go:57 | 🟢 低 | 短信验证码 MOCK 到控制台，可接受但需确保短信服务接口预留清晰 | 已有 MOCK 标注，合理 |
| 6 | 管理端 | pages/Dashboard.tsx | 🟡 中 | 趋势图表区域为纯文字占位 `图表区域（集成图表库后展示...）`，未集成图表库 | 集成 @ant-design/charts 或 echarts，至少展示静态示意数据 |
| 7 | 管理端 | pages/data/Board.tsx | 🟡 中 | 数据看板页依赖 `getDashboardCharts` API，但后端未实现该路由 | 后端增加 `/admin/dashboard/charts` 接口或管理端做容错处理 |

---

## 逐模块分析

### 后端 (backend/) ✅ 整体良好

**编译 & 测试:** `go build` + `go test` 全部通过  
**分层规范:** 严格遵守 api→service→repository 分层  
**安全:**
- ✅ 密码 bcrypt 加密
- ✅ 手机号脱敏（138****8000）
- ✅ JWT认证完整
- ✅ SQL参数化查询（GORM）

**API对齐:** 57个路由全部注册且有handler实现，仅缺少3个PRD定义的管理后台路由（问题1-3）

**代码中的 `return nil`:** 检查后为正常的成功返回（ReplaceTags/ReplaceResources 操作完成后返回nil表示无错误），非空实现。

### 前端 (frontend/) ✅ 完整

**构建:** pnpm build 零错误  
**空实现/TODO:** 无（grep搜索仅命中input的placeholder属性，正常）  
**API对接:** services.ts 覆盖全部后端公开+认证接口  
**按钮功能:** 所有按钮有onClick实现（购买、收藏、点赞、搜索、登录、注册、翻页等）  
**响应式:** Tailwind grid-cols-2/md:3/lg:4 完整  
**错误处理:** 401拦截跳转登录、catch块处理异常  

### 管理端 (admin/) ✅ 整体良好

**构建:** pnpm build 通过  
**空实现/TODO:** 无  
**API对接:** 40+ API调用全部实现，注释标明路径  
**表单校验:** antd Form.Item rules  
**防重复提交:** loading状态控制  
**二次确认:** Modal.confirm  
**问题:** 图表占位（问题6）+ charts API缺失（问题7）

---

## 前后端对齐检查

| 检查项 | 状态 |
|--------|------|
| 前端API baseURL 与后端端口一致 | ✅ localhost:8080/api/v1 |
| 响应格式 {code, message, data} | ✅ 一致 |
| 分页响应 {items, total, page, pageSize, totalPages} | ✅ 一致 |
| 错误码 40001/40101 | ✅ 一致 |
| 字段命名 snake_case | ✅ 一致（course_id, is_vip, captcha_key等） |
| Token 传递 Authorization: Bearer | ✅ 一致 |
| 管理端API路径前缀 /admin/ | ✅ 一致 |

---

## 结论

核心功能代码完整度高，无阻塞性空实现。需修复的7个问题中：
- **3个后端路由缺失**（高/中优先级）
- **1个MOCK需完善返回格式**
- **2个管理端图表需补充**
- **1个可接受的MOCK**

建议指挥官下发修复指令。
