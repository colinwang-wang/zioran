# 知猿 (Zioran) — 项目上下文

## 项目概述
网课资源付费下载站。用户注册→充值金币→购买课程/VIP→下载资源。

- **域名**: zioran.com
- **服务器**: 120.26.192.163 (Ubuntu, SSH root/Admin66668888)
- **Git**: git@github.com:colinwang-wang/zioran.git
- **在线访问**: http://120.26.192.163 (前端) | http://120.26.192.163/admin/ (管理后台 admin/admin123456)

## 技术栈

| 层 | 技术 |
|----|------|
| 前端用户端 | Next.js 14 (App Router) + Tailwind CSS |
| 管理后台 | React 18 + Vite + Ant Design 5 (base: /admin/) |
| 后端 API | Go + Gin + GORM |
| 数据库 | MySQL 8 (root/root123456, 库名 zioran) |
| 缓存 | Redis |
| 进程管理 | pm2 (zioran-api + zioran-frontend) |
| 反向代理 | Nginx (80 → 3000前端, /api/ → 8080后端, /admin/ → 静态) |

## 目录结构

```
zioran/
├── backend/                    # Go 后端
│   ├── cmd/server/main.go      # 入口
│   ├── internal/
│   │   ├── api/                # Handler + Router (77个路由)
│   │   ├── service/            # 业务逻辑
│   │   ├── repository/         # 数据访问
│   │   ├── model/              # 数据模型 + DTO
│   │   └── middleware/         # JWT
│   ├── pkg/
│   │   ├── sms/                # 短信服务(mock/aliyun/tencent)
│   │   ├── payment/            # 微信支付+支付宝(config切换)
│   │   ├── oauth/              # 微信登录OAuth
│   │   ├── errcode/            # 统一错误码
│   │   ├── response/           # 统一响应 {code,message,data}
│   │   └── config/             # 配置加载
│   ├── migrations/             # MySQL建表SQL(001-004)
│   ├── config.yaml             # 配置文件
│   └── go.mod
├── frontend/                   # Next.js 用户端
│   ├── src/app/                # 页面路由(19个)
│   ├── src/components/         # 公共组件
│   ├── src/lib/api.ts          # Axios封装(SSR用127.0.0.1,客户端用/api/v1)
│   ├── src/lib/services.ts     # API方法
│   ├── src/types/index.ts      # TypeScript类型
│   └── src/contexts/AuthContext.tsx
├── admin/                      # React 管理后台
│   ├── src/pages/              # 14个管理页面
│   ├── src/api/index.ts        # API封装
│   ├── src/layouts/            # ProLayout
│   └── vite.config.ts          # base: '/admin/'
├── docs/
│   ├── prd/                    # PRD文档(00-09)
│   ├── prototype/              # HTML原型(8页)
│   ├── DESIGN.md               # 设计系统(Pinterest风格,主色#ff0036)
│   └── fix-list.md             # 调整清单
├── .skills/                    # 开发规范(8个skill)
├── .commander/                 # 多专家协作协议
├── Makefile                    # 运维命令
└── CLAUDE.md                   # 项目入口
```

## 核心业务流程

```
用户注册(手机号+短信) → 登录 → 浏览课程 → 充值金币(微信/支付宝)
  → 购买课程(扣金币) → 下载资源
  → 或购买VIP(99金币) → 全站免费下载
```

## API 响应格式
```json
{"code": 0, "message": "ok", "data": {...}}
// 分页: data = {items:[], total, page, pageSize, totalPages}
// 错误码: 40001参数错误, 40101未认证, 40301无权限, 40401不存在, 50001服务器错误
```

## 关键设计决策

1. **前端API解包**: api.ts拦截器自动从`{code,data}`中提取data，services.ts直接拿到业务数据
2. **SSR数据获取**: page.tsx用原生fetch(http://127.0.0.1:8080)，需手动`.json().then(r=>r.data)`提取
3. **验证码**: 生成SVG图片返回`data:image/svg+xml;base64,...`，测试万能码`0000`
4. **短信**: 测试万能码`000000`跳过验证，真实短信通过config.yaml配置provider切换
5. **支付**: config中`enabled:false`时走MOCK(直接完成)，填入密钥后真实支付
6. **管理员登录**: 独立接口`POST /admin/login`(username+password)，与用户登录(phone+password+captcha)分离

## 外部服务配置 (config.yaml)

```yaml
sms.provider: mock          # → aliyun | tencent (填密钥即生效)
payment.wechat.enabled: false   # → true (需AppID/MchID/APIKey)
payment.alipay.enabled: false   # → true (需AppID/PrivateKey)
oauth.wechat.enabled: false     # → true (需AppID/AppSecret)
```

## 运维命令 (Makefile)

```bash
make deploy-full    # 一键全量部署
make quick-deploy   # 快速部署(仅后端)
make status         # pm2服务状态
make logs-api       # 后端日志
make restart        # 重启
make health         # 健康检查
make db-backup      # 数据库备份
make ssh            # SSH连接服务器
```

## 已知待处理

- 域名 zioran.com 需ICP备案后解析到120.26.192.163
- 微信/支付宝/短信需企业资质申请后填入config.yaml
- 课程数据需通过管理后台录入真实内容
- 当前有测试bypass(验证码0000/短信000000)，上线前需移除或改为环境变量控制

## 数据库种子数据

- 14个课程分类
- 1个VIP套餐(99金币/原价699/永久)
- 12个示例课程(含封面图)
- 1个管理员(admin/admin123456)
- 6个金刚区导航
- 2个Banner
