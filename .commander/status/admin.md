# Admin 管理后台 - 状态报告

## 状态: ✅ 完成

## 构建验证
- [x] `pnpm build` 零错误通过
- [x] TypeScript 类型检查通过

## 技术栈
- React 18 + TypeScript + Vite
- Ant Design 5 + ProLayout
- React Router 6 (HashRouter)
- Axios 封装 (JWT Bearer, baseURL=/api/v1)
- 路由懒加载 (所有页面组件)

## 页面清单

| 页面 | 路径 | 功能 |
|------|------|------|
| 登录 | /login | 用户名+密码登录, 表单校验 |
| 仪表盘 | /dashboard | 4统计卡片 + 趋势图区域 + 最近订单表 |
| 课程列表 | /courses | 搜索/分类/状态筛选 + 分页 + 批量上下架/删除 |
| 课程编辑 | /courses/create, /courses/:id/edit | 全字段表单(主图/标题/分类/标签/价格/详情/资源/状态) |
| 分类管理 | /categories | 列表 + Modal新增/编辑 + 上下架 |
| 标签管理 | /tags | 列表 + Modal CRUD |
| 用户列表 | /users | VIP/普通筛选 + 手动充值 + 禁用/启用 |
| 用户详情 | /users/:id | 信息卡片展示(脱敏手机号) |
| 订单列表 | /orders | 类型/状态/时间范围筛选 + 分页 |
| 订单详情 | /orders/:id | 订单信息 + 退款操作 |
| 留言管理 | /guestbook | 隐藏/显示 + 置顶/取消 + 删除 |
| 评论管理 | /comments | 隐藏/显示 + 删除 |
| 首页配置 | /config | 金刚区CRUD + Banner CRUD (含图片上传) |
| 数据看板 | /data | 统计卡片 + 柱状图(日/月/季/年) |

## 规范遵循
- [x] 搜索条件变更重置页码
- [x] 表单防重复提交 (loading状态控制按钮)
- [x] 删除/危险操作二次确认 (Modal.confirm)
- [x] 手机号脱敏展示
- [x] API路径从PRD复制, 注释标明完整路径
- [x] 无空实现/TODO占位
- [x] 401自动跳转登录页

## 文件结构
```
admin/
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
└── src/
    ├── main.tsx
    ├── App.tsx
    ├── vite-env.d.ts
    ├── api/index.ts
    ├── types/index.ts
    ├── utils/request.ts
    ├── layouts/AdminLayout.tsx
    └── pages/
        ├── Login.tsx
        ├── Dashboard.tsx
        ├── course/{List,Form}.tsx
        ├── category/List.tsx
        ├── tag/List.tsx
        ├── user/{List,Detail}.tsx
        ├── order/{List,Detail}.tsx
        ├── guestbook/List.tsx
        ├── comment/List.tsx
        ├── config/HomeConfig.tsx
        └── data/Board.tsx
```
