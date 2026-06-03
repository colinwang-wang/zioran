# Backend Status - Phase 3+4

## 状态: ✅ 完成

## 完成日期: 2026-06-04

## 交付清单

### 1. 金币接口 ✅
- `GET /api/v1/coins/balance` - 查询余额
- `GET /api/v1/coins/transactions` - 交易记录(分页)
- `POST /api/v1/coins/recharge` - 充值(MOCK支付回调,立即到账)

### 2. VIP接口 ✅
- `GET /api/v1/vip/packages` - 套餐列表(公开)
- `GET /api/v1/vip/status` - VIP状态
- `POST /api/v1/vip/purchase` - 扣金币开通VIP

### 3. 订单接口 ✅
- `POST /api/v1/orders` - 创建订单(course类型)
- `GET /api/v1/orders/:id` - 订单详情
- `POST /api/v1/courses/:id/download` - 获取下载链接(需已购买或VIP)

### 4. 个人中心 ✅
- `GET /api/v1/user/orders` - 我的订单
- `GET /api/v1/user/downloads` - 我的下载
- `GET /api/v1/user/favorites` - 我的收藏
- `POST /api/v1/user/favorites` - 添加收藏
- `DELETE /api/v1/user/favorites/:courseId` - 取消收藏
- `PUT /api/v1/user/password` - 修改密码

### 5. 留言板 ✅
- `GET /api/v1/guestbook` - 留言列表(公开)
- `POST /api/v1/guestbook` - 发布留言
- `POST /api/v1/guestbook/:id/like` - 点赞/取消
- `DELETE /api/v1/guestbook/:id` - 删除自己的留言

### 6. 评论 ✅
- `GET /api/v1/comments?targetType=course&targetId=xxx` - 评论列表
- `POST /api/v1/comments` - 发表评论(支持嵌套回复)
- `DELETE /api/v1/comments/:id` - 删除评论

### 7. 首页配置 ✅
- `GET /api/v1/home/nav-items` - 金刚区导航(公开)
- `GET /api/v1/home/banners` - Banner列表(公开)

### 8. 后台管理 ✅
- 订单管理: `GET /api/v1/admin/orders`, `POST /api/v1/admin/orders/:id/refund`
- 用户管理: `GET /api/v1/admin/users`, `PUT status`, `POST recharge`
- 留言管理: `GET/PUT status/PUT pin/DELETE`
- 评论管理: `GET/PUT status/DELETE`
- 首页配置: nav-items CRUD, banners CRUD
- 仪表盘: `GET /api/v1/admin/dashboard/stats`

## 核心流程验证

### 购买课程全流程 ✅
1. 充值金币 → 余额增加
2. 创建订单(type=course) → 扣减金币 + 创建purchase记录
3. 获取下载链接 → 返回资源URL+提取码
4. 金币不会为负数(事务保护)

### VIP购买流程 ✅
1. 查看套餐 → 选择终身VIP(99金币)
2. 扣金币 → 开通VIP → 状态变更
3. VIP用户免费下载所有课程

## 技术实现

- 金币扣减使用数据库事务 + FOR UPDATE 行锁
- 余额检查在事务内完成,不会出现负数
- MOCK支付: 充值后立即回调成功
- 分层架构: handler → service → repository

## 新增文件
- `internal/model/payment.go` - 金币/VIP/订单/留言/评论/首页模型
- `internal/model/payment_dto.go` - 请求/响应DTO
- `internal/repository/payment.go` - 支付相关数据访问
- `internal/repository/community.go` - 社区/首页配置数据访问
- `internal/service/payment.go` - 支付业务逻辑
- `internal/service/community.go` - 社区业务逻辑
- `internal/api/payment.go` - 支付API处理器
- `internal/api/community.go` - 社区API处理器
- `internal/api/admin_payment.go` - 后台管理API处理器
- `internal/api/payment_test.go` - 22个测试用例
- `migrations/002_phase3_4.sql` - 新增表结构

## 测试结果
```
ok  github.com/zioran/backend/internal/api  2.550s
```
全部测试通过，包括 Phase 1/2 的回归测试。
