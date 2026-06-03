# Phase 1 — 后端专家指令

> 状态: PENDING
> 依赖: 无
> 更新时间: 2026-06-03T23:41

## 背景
知猿(zioran)项目启动，需要搭建 Go+Gin 后端基础架构，完成数据库建表和用户认证模块。

## 项目信息
- 数据库：MySQL 8, root/root123456, 数据库名 zioran
- 技术栈：Go + Gin + GORM + Redis + JWT
- 参考文档：docs/prd/01-user.md, docs/prd/07-database.md, docs/prd/08-api.md

## 任务列表

### 1. 项目初始化
在项目根目录创建 `backend/` 目录，初始化 Go module：
```
backend/
├── cmd/server/main.go
├── internal/
│   ├── api/          # 路由 + handler
│   ├── service/      # 业务逻辑
│   ├── repository/   # 数据访问
│   ├── model/        # 数据库模型 + DTO
│   └── middleware/   # JWT、CORS
├── pkg/
│   ├── errcode/      # 统一错误码
│   ├── response/     # 统一响应 {code, message, data}
│   └── config/       # 配置加载
├── migrations/       # 数据库迁移 SQL
├── config.yaml       # 配置文件
└── go.mod
```

### 2. 数据库迁移
参考 docs/prd/07-database.md，将 PostgreSQL 语法转为 MySQL 语法，创建以下核心表：
- users（手机号为主标识，NOT NULL UNIQUE）
- categories
- tags
- courses + course_tags + course_resources
- coin_accounts + coin_transactions
- vip_packages + user_vip
- orders + purchases

注意 MySQL 语法差异：BIGSERIAL → BIGINT AUTO_INCREMENT, TIMESTAMP → DATETIME, TEXT 保持, JSONB → JSON, BOOLEAN → TINYINT(1)

### 3. 用户认证模块（TDD）
严格按 RED-GREEN-REFACTOR 开发以下接口：

| 接口 | 说明 |
|------|------|
| POST /api/v1/auth/captcha | 获取图形验证码 |
| POST /api/v1/auth/sms/send | 发送短信验证码（暂用 mock，控制台打印） |
| POST /api/v1/auth/register | 手机号+短信验证码+密码注册 |
| POST /api/v1/auth/login | 手机号+密码+图形验证码登录 |
| GET /api/v1/user/profile | 获取当前用户信息（需JWT） |

### 4. 统一响应格式
```json
{ "code": 0, "message": "ok", "data": {} }
```
错误码：40001 参数错误, 40101 未认证, 40301 无权限, 50001 服务器错误

### 5. 配置文件 config.yaml
```yaml
server:
  port: 8080
database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: root123456
  dbname: zioran
redis:
  host: 127.0.0.1
  port: 6379
jwt:
  secret: zioran-jwt-secret-2026
  expire: 72h
```

## 交付标准
- [ ] `go build ./...` 零错误
- [ ] `go test ./...` 全部通过（认证模块覆盖率 ≥ 80%）
- [ ] 数据库表已创建（通过 migration 文件）
- [ ] 注册→登录→获取profile 主流程可用
- [ ] 短信验证码 mock 到控制台输出
- [ ] 统一响应格式，错误返回正确业务码

## 自检清单
- [ ] 手机号重复注册返回 40001
- [ ] 密码错误返回 40001
- [ ] 缺少必填字段返回 40001
- [ ] 未携带 Token 访问 profile 返回 40101
- [ ] Token 正确时返回用户信息（手机号脱敏）
