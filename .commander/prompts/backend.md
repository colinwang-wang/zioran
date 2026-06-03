# Phase 6 — 后端专家指令：补齐PRD定义的剩余接口

> 状态: PENDING
> 更新时间: 2026-06-04T01:05
> 背景: QA对比PRD发现25个定义但未实现的API，按优先级分批补齐

## 优先级 A — 核心功能（必须实现）

### 1. 忘记密码
```
POST /api/v1/auth/forgot-password
Body: { phone, sms_code, new_password }
→ 验证短信码 → 更新密码
```

### 2. Token刷新
```
POST /api/v1/auth/refresh
Header: Authorization: Bearer <token>
→ 验证旧token → 生成新token
```

### 3. 订单取消
```
POST /api/v1/orders/:id/cancel
→ 仅pending状态可取消
```

### 4. 用户订单详情
```
GET /api/v1/user/orders/:id
→ 返回订单详情（仅本人可查看）
```

### 5. 工单系统（用户端）
```
GET  /api/v1/tickets            — 我的工单列表
POST /api/v1/tickets            — 提交工单 { title, content }
GET  /api/v1/tickets/:id        — 工单详情（含回复）
POST /api/v1/tickets/:id/reply  — 回复工单 { content }
```

### 6. 工单系统（管理端）
```
GET  /api/v1/admin/tickets             — 全部工单列表
GET  /api/v1/admin/tickets/:id         — 工单详情
PUT  /api/v1/admin/tickets/:id/status  — 更新状态(processing/replied/closed)
POST /api/v1/admin/tickets/:id/reply   — 管理员回复
```

## 优先级 B — 管理后台辅助功能

### 7. 系统设置
```
GET /api/v1/admin/settings   — 获取全部设置
PUT /api/v1/admin/settings   — 更新设置 { key: value }
```

### 8. 管理员账号管理
```
GET    /api/v1/admin/admins      — 管理员列表
POST   /api/v1/admin/admins      — 创建管理员 { username, password, role }
PUT    /api/v1/admin/admins/:id  — 更新管理员
DELETE /api/v1/admin/admins/:id  — 删除管理员
```

### 9. 财务管理
```
GET /api/v1/admin/finance/summary      — 收支明细（每日/已结算/待结算）
GET /api/v1/admin/finance/withdrawals  — 提现列表
```

### 10. 日志
```
GET /api/v1/admin/logs/operations  — 管理员操作日志
GET /api/v1/admin/logs/payments    — 支付异常日志
```

### 11. 评论管理员回复
```
POST /api/v1/admin/comments/:id/reply  — 管理员回复评论 { content }
```

## 优先级 C — 第三方对接预留（MOCK）

### 12. 支付回调
```
POST /api/v1/pay/notify/wechat   — 微信回调（MOCK: 直接标记订单已支付）
POST /api/v1/pay/notify/alipay   — 支付宝回调（MOCK: 同上）
```

### 13. 微信登录
```
GET  /api/v1/auth/oauth/wechat           — 返回微信授权URL（MOCK）
POST /api/v1/auth/oauth/wechat/callback  — 微信回调（MOCK: 自动创建用户）
```

### 14. 其他
```
POST /api/v1/upload/images    — 批量上传（复用单图上传逻辑）
GET  /api/v1/home/config      — 首页配置（文案等，从settings表读取）
```

## 交付标准
- [ ] go build ./... 零错误
- [ ] go test ./... 全部通过
- [ ] 所有新路由在 router.go 注册
- [ ] MOCK接口标注 `// MOCK: 待接入真实服务`
- [ ] 工单表 tickets + ticket_replies 通过migration创建

## 同时需要
- 管理端补充工单管理页面
- 管理端补充系统设置页面
- 前端补充忘记密码页面、工单页面
