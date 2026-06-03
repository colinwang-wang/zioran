# 08 - API 设计

## 8.1 通用规范

### Base URL
```
开发环境: http://localhost:8080/api/v1
生产环境: https://api.zioran.com/api/v1
```

### 认证
- Header: `Authorization: Bearer <jwt_token>`
- 公开接口不需要

### 通用响应格式

成功：
```json
{
  "code": 0,
  "message": "ok",
  "data": { ... }
}
```

分页列表：
```json
{
  "code": 0,
  "data": {
    "items": [...],
    "total": 3200,
    "page": 1,
    "pageSize": 16,
    "totalPages": 200
  }
}
```

错误：
```json
{
  "code": 40001,
  "message": "用户名已存在"
}
```

### 错误码
| 范围 | 说明 |
|------|------|
| 40001-40099 | 参数校验错误 |
| 40100-40199 | 认证错误 |
| 40300-40399 | 权限错误 |
| 40400-40499 | 资源不存在 |
| 50000-50099 | 服务器错误 |

---

## 8.2 认证接口

```
POST   /api/v1/auth/captcha           # 获取图形验证码
POST   /api/v1/auth/sms/send          # 发送短信验证码（需先验图形验证码）
POST   /api/v1/auth/register          # 手机号注册（手机号+短信验证码+密码）
POST   /api/v1/auth/login             # 手机号+密码登录
POST   /api/v1/auth/forgot-password   # 忘记密码（手机号+短信验证码+新密码）
POST   /api/v1/auth/refresh           # 刷新token
GET    /api/v1/auth/oauth/wechat      # 微信登录授权URL
POST   /api/v1/auth/oauth/wechat/callback  # 微信登录回调
```

---

## 8.3 用户接口

```
GET    /api/v1/user/profile            # 个人信息
PUT    /api/v1/user/profile            # 更新信息
PUT    /api/v1/user/password           # 修改密码
GET    /api/v1/user/orders             # 我的订单
GET    /api/v1/user/orders/:id         # 订单详情
GET    /api/v1/user/downloads          # 我的下载
GET    /api/v1/user/favorites          # 我的收藏
POST   /api/v1/user/favorites          # 添加收藏
DELETE /api/v1/user/favorites/:courseId # 取消收藏
```

---

## 8.4 课程接口

```
GET    /api/v1/courses                 # 课程列表
       Query: page, pageSize, categoryId, tagId, sort, keyword
GET    /api/v1/courses/:slug           # 课程详情
POST   /api/v1/courses/:id/like        # 收藏/取消
POST   /api/v1/courses/:id/download    # 获取下载链接
GET    /api/v1/courses/latest          # 最新课程(limit=8)
GET    /api/v1/categories              # 分类列表
GET    /api/v1/tags                    # 标签列表
```

---

## 8.5 搜索接口

```
GET    /api/v1/search
       Query: q, categoryId, page, pageSize
```

---

## 8.6 金币接口

```
GET    /api/v1/coins/balance           # 余额
GET    /api/v1/coins/transactions      # 交易记录
POST   /api/v1/coins/recharge          # 创建充值订单
       Body: { amount, payMethod }
```

---

## 8.7 VIP 接口

```
GET    /api/v1/vip/packages            # 套餐列表
GET    /api/v1/vip/status              # VIP状态
POST   /api/v1/vip/purchase            # 购买VIP
       Body: { packageId }
```

---

## 8.8 订单接口

```
POST   /api/v1/orders                  # 创建订单
GET    /api/v1/orders/:id              # 订单详情
POST   /api/v1/orders/:id/cancel       # 取消订单
```

---

## 8.9 支付回调

```
POST   /api/v1/pay/notify/wechat       # 微信支付回调
POST   /api/v1/pay/notify/alipay       # 支付宝回调
```

---

## 8.10 留言板接口

```
GET    /api/v1/guestbook               # 留言列表
POST   /api/v1/guestbook               # 发布留言
POST   /api/v1/guestbook/:id/like      # 点赞/取消
DELETE /api/v1/guestbook/:id           # 删除自己的留言
```

---

## 8.11 评论接口

```
GET    /api/v1/comments                # 评论列表
       Query: targetType, targetId, page, pageSize
POST   /api/v1/comments                # 发表评论
DELETE /api/v1/comments/:id            # 删除自己的评论
```

---

## 8.12 工单接口

```
GET    /api/v1/tickets                 # 我的工单列表
POST   /api/v1/tickets                 # 提交工单
       Body: { title, content }
GET    /api/v1/tickets/:id             # 工单详情（含回复）
POST   /api/v1/tickets/:id/reply       # 回复工单
       Body: { content }
```

---

## 8.13 首页配置接口（公开）

```
GET    /api/v1/home/nav-items          # 金刚区导航
GET    /api/v1/home/banners            # Banner列表
GET    /api/v1/home/config             # 首页配置（文案等）
```

---

## 8.14 管理后台接口

```
# 课程管理
GET    /api/v1/admin/courses
POST   /api/v1/admin/courses
PUT    /api/v1/admin/courses/:id
DELETE /api/v1/admin/courses/:id
PUT    /api/v1/admin/courses/:id/status
POST   /api/v1/admin/courses/batch

# 分类管理
GET    /api/v1/admin/categories
POST   /api/v1/admin/categories
PUT    /api/v1/admin/categories/:id
DELETE /api/v1/admin/categories/:id
PUT    /api/v1/admin/categories/:id/status  # 上下架

# 标签管理
GET    /api/v1/admin/tags
POST   /api/v1/admin/tags
PUT    /api/v1/admin/tags/:id
DELETE /api/v1/admin/tags/:id

# 订单管理
GET    /api/v1/admin/orders
GET    /api/v1/admin/orders/:id
POST   /api/v1/admin/orders/:id/refund

# 用户管理
GET    /api/v1/admin/users
GET    /api/v1/admin/users/:id
PUT    /api/v1/admin/users/:id
PUT    /api/v1/admin/users/:id/status      # 禁用/启用
POST   /api/v1/admin/users/:id/recharge    # 手动充值

# 留言管理
GET    /api/v1/admin/guestbook
PUT    /api/v1/admin/guestbook/:id/status
PUT    /api/v1/admin/guestbook/:id/pin
DELETE /api/v1/admin/guestbook/:id

# 评论管理
GET    /api/v1/admin/comments
PUT    /api/v1/admin/comments/:id/status
POST   /api/v1/admin/comments/:id/reply    # 管理员回复
DELETE /api/v1/admin/comments/:id

# 工单管理
GET    /api/v1/admin/tickets
GET    /api/v1/admin/tickets/:id
PUT    /api/v1/admin/tickets/:id/status
POST   /api/v1/admin/tickets/:id/reply

# 财务管理
GET    /api/v1/admin/finance/summary       # 收支明细
GET    /api/v1/admin/finance/withdrawals   # 提现列表

# 数据看板
GET    /api/v1/admin/dashboard/stats       # 统计数据
GET    /api/v1/admin/dashboard/charts      # 图表数据
       Query: period (day/month/quarter/year)

# 系统配置
GET    /api/v1/admin/settings
PUT    /api/v1/admin/settings
GET    /api/v1/admin/nav-items
POST   /api/v1/admin/nav-items
PUT    /api/v1/admin/nav-items/:id
DELETE /api/v1/admin/nav-items/:id
GET    /api/v1/admin/banners
POST   /api/v1/admin/banners
PUT    /api/v1/admin/banners/:id
DELETE /api/v1/admin/banners/:id

# 权限管理
GET    /api/v1/admin/admins
POST   /api/v1/admin/admins
PUT    /api/v1/admin/admins/:id
DELETE /api/v1/admin/admins/:id

# 日志
GET    /api/v1/admin/logs/operations       # 操作日志
GET    /api/v1/admin/logs/payments         # 支付异常日志
```

---

## 8.15 文件上传

```
POST   /api/v1/upload/image            # 上传图片
       限制: 最大 10MB, jpg/png/webp/gif
POST   /api/v1/upload/images           # 批量上传
       限制: 最多 20 张
```
