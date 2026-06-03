# 01 - 用户模块

## 1.1 登录页

### 路由
`/login`

### 页面结构
```
┌──────────────────────────────────────────┐
│  div.sign (弹窗模式，全屏居中)             │
│  ┌──────────────────────────────────┐    │
│  │  div.sign-mask (半透明遮罩)       │    │
│  │  div.sign-box                    │    │
│  │  ┌────────────────────────────┐  │    │
│  │  │ Logo 图片                   │  │    │
│  │  ├────────────────────────────┤  │    │
│  │  │ [手机号]      icon: phone   │  │    │
│  │  ├────────────────────────────┤  │    │
│  │  │ [密码]        icon: lock    │  │    │
│  │  ├────────────────────────────┤  │    │
│  │  │ [图形验证码] [显码] icon:safe│  │    │
│  │  ├────────────────────────────┤  │    │
│  │  │ [登录] 按钮                 │  │    │
│  │  ├────────────────────────────┤  │    │
│  │  │ 没有账号？注册              │  │    │
│  │  │ 忘记密码？                  │  │    │
│  │  ├────────────────────────────┤  │    │
│  │  │ 社交账号快速登录             │  │    │
│  │  │ [微信] 图标                 │  │    │
│  │  └────────────────────────────┘  │    │
│  └──────────────────────────────────┘    │
└──────────────────────────────────────────┘
```

### 表单字段
| 字段 | name | 类型 | 必填 | placeholder |
|------|------|------|------|-------------|
| 手机号 | phone | tel | 是 | 手机号 |
| 密码 | password | password | 是 | 密码 |
| 图形验证码 | captcha | text | 是 | 验证码 |

### 验证码
- 初始显示文字 "显示验证码"
- 点击后显示 4 位图片验证码
- 可点击刷新

### 交互
- 点击"登录" → POST 提交 → 成功关闭弹窗刷新页面 / 失败显示错误提示
- "注册" → 切换到注册表单（同一弹窗内切换）
- "忘记密码？" → 跳转 `/login?action=password`
- "微信" 图标 → 微信快捷登录（扫码或授权）

### API
```
POST /api/v1/auth/login
Content-Type: application/json
Body: {
  "phone": "string",
  "password": "string",
  "captcha": "string",
  "captcha_key": "string"
}
Response 200: {
  "code": 0,
  "data": {
    "token": "jwt_string",
    "user": { "id": 1, "username": "xxx", "phone": "138****8888", "avatar": "url", "is_vip": false }
  }
}
Response 400: { "code": 40001, "message": "手机号或密码错误" }
```

---

## 1.2 注册页

### 路由
`/register`

### 页面结构（弹窗内切换到注册表单）
```
form#sign-up
├── Logo 图片
├── [手机号]       icon: phone
├── [短信验证码] [获取验证码]  icon: message
├── [密码]       icon: lock
├── [注册] 按钮
├── 返回登录
└── 社交账号快速登录
    └── [微信]
```

### 表单字段
| 字段 | name | 类型 | 必填 | placeholder |
|------|------|------|------|-------------|
| 手机号 | phone | tel | 是 | 手机号 |
| 短信验证码 | sms_code | text | 是 | 短信验证码 |
| 密码 | password | password | 是 | 设置密码 |

### 短信验证码
- 点击"获取验证码"→发送短信→按钮变为60秒倒计时
- 60秒内不可重复发送
- 验证码5分钟有效

### API
```
POST /api/v1/auth/sms/send
Body: { "phone": "string", "captcha": "string", "captcha_key": "string" }
Response: { "code": 0, "message": "验证码已发送" }

POST /api/v1/auth/register
Body: {
  "phone": "string",
  "sms_code": "string",
  "password": "string"
}
Response 200: { "code": 0, "data": { "token": "jwt_string", "user": {...} } }
Response 400: { "code": 40001, "message": "手机号已注册" }
```

---

## 1.3 验证码接口

```
GET /api/v1/auth/captcha
Response: {
  "code": 0,
  "data": {
    "captcha_key": "uuid_string",
    "captcha_image": "base64_png_string"
  }
}
```
- 验证码 5 分钟有效
- 每次请求生成新的 key

---

## 1.4 忘记密码

### 路由
`/login?action=password`

### 页面文案
**求学选知猿，精品网课伴成长**

### 流程
1. 输入注册手机号
2. 获取短信验证码
3. 输入验证码 + 新密码 → 重置

### API
```
POST /api/v1/auth/forgot-password
Body: { "phone": "string", "sms_code": "string", "new_password": "string" }
Response: { "code": 0, "message": "密码已重置" }
```

---

## 1.5 快捷登录 - 微信

### 流程
1. 用户点击微信图标
2. 展示微信登录二维码（或跳转微信授权页）
3. 用户扫码/确认授权
4. 回调获取用户微信信息
5. 关联账号或自动注册

### API
```
GET  /api/v1/auth/oauth/wechat?redirect_url=当前页
     → 返回微信授权 URL

POST /api/v1/auth/oauth/wechat/callback
Body: { "code": "string" }
Response: { "code": 0, "data": { "token": "jwt_string", "user": {...}, "is_new": true/false } }
```

---

## 1.6 登录态管理

### 方式
- 独立页面 `/login` + `/register`（也可做弹窗）
- JWT token 存储在 localStorage
- Axios 拦截器自动附加 `Authorization: Bearer <token>`
- 前端全局状态管理用户信息

### 用户状态 Hook
```typescript
interface AuthState {
  isLoggedIn: boolean;
  user: User | null;
  token: string | null;
}

interface User {
  id: number;
  username: string;
  phone: string;
  avatar: string;
  is_vip: boolean;
  coin_balance: number;
}
```

---

## 1.7 个人中心

### 路由
`/user`

### 参考
- https://chenwenb.com/（简版，以简版为主）
- https://www.zycku.com/（复杂版）

### 子页面
| 路由 | 页面 | 说明 |
|------|------|------|
| `/user` | 概览 | 头像、用户名、VIP状态、金币余额、最近订单 |
| `/user/orders` | 我的订单 | 购买记录列表 |
| `/user/downloads` | 我的下载 | 已购资源列表 |
| `/user/favorites` | 我的收藏 | 收藏课程列表 |
| `/user/recharge` | 充值 | 微信和支付宝充值 |
| `/user/settings` | 账号设置 | 修改密码、邮箱 |

> 注：不需要"我的推广"功能，其他参考网站功能都保留。

### 充值
- 支持微信支付和支付宝充值
- 充值页面展示充值档位和支付方式选择

### API
```
GET    /api/v1/user/profile
PUT    /api/v1/user/profile
PUT    /api/v1/user/password
GET    /api/v1/user/orders?page=1&pageSize=10
GET    /api/v1/user/orders/:id
GET    /api/v1/user/downloads?page=1&pageSize=10
GET    /api/v1/user/favorites?page=1&pageSize=10
POST   /api/v1/user/favorites
DELETE /api/v1/user/favorites/:courseId
```

---

## 1.8 权限体系

### 角色
| 角色 | 权限 |
|------|------|
| visitor | 浏览课程列表、详情（不含下载区） |
| user | 浏览、购买、下载、留言、评论、收藏 |
| vip | 全站免费下载（注册成为会员，前端价格显示为 0） |
| admin | 全部权限 |

### 课程详情页权限状态机
```
visitor:
  下载区 → "请先登录" + 登录链接
  评论区 → "请先登录"

user (非VIP, 未购买):
  下载区 → "下载价格: X 金币" + "终身VIP免费" + "立即购买"按钮
  评论区 → 可评论

user (VIP):
  下载区 → "免费下载" + 下载链接（价格显示为 0）

user (非VIP, 已购买):
  下载区 → "已购买" + 下载链接
```
