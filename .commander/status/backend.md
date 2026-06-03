# 后端修复状态报告

> 时间: 2026-06-04T00:31
> 状态: ✅ PASS

## 修复内容

| # | 问题 | 修复方式 | 状态 |
|---|------|----------|------|
| 1 | `GET /admin/users/:id` 缺失 | 添加 `AdminGetUser` service + `UserDetail` handler + 路由注册 | ✅ |
| 2 | `PUT /admin/categories/:id/status` 缺失 | 添加 `AdminCategoryUpdateStatus` service + handler + 路由注册 | ✅ |
| 3 | `POST /admin/courses/batch` 缺失 | 添加 `AdminCourseBatch` service + `CourseBatch` handler + 路由注册 | ✅ |
| 4 | 充值MOCK直接成功，缺少pay_url | 修改返回 `{order_id, order_no, pay_url: "mock://pay"}`，MOCK仍完成支付 | ✅ |
| 5 | `GET /admin/dashboard/charts` 缺失 | 添加 `DashboardCharts` service + handler + 路由注册，返回示意数据 | ✅ |

## 修改文件

- `internal/api/router.go` — 新增4条路由
- `internal/api/admin.go` — 新增 `CategoryUpdateStatus`、`CourseBatch` handler
- `internal/api/admin_payment.go` — 新增 `UserDetail`、`DashboardCharts` handler
- `internal/service/payment.go` — 修改 `Recharge` 返回格式，新增 `AdminGetUser`、`DashboardCharts`
- `internal/service/course.go` — 新增 `AdminCategoryUpdateStatus`、`AdminCourseBatch`
- `internal/model/payment_dto.go` — 新增 `RechargeResponse`、`ChartPoint`、`DashboardChartsResponse`、`AdminCourseBatchRequest`、`AdminCategoryStatusRequest`

## 验证

- `go build ./...` ✅ 零错误
- `go test ./...` ✅ 全部通过
- 新路由已在 router.go 注册

## 新增路由清单

```
GET    /api/v1/admin/users/:id              → UserDetail
PUT    /api/v1/admin/categories/:id/status  → CategoryUpdateStatus
POST   /api/v1/admin/courses/batch          → CourseBatch
GET    /api/v1/admin/dashboard/charts       → DashboardCharts
```
