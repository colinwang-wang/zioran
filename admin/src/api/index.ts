import request from '@/utils/request'
import type { ApiResponse, PageResponse, LoginParams, LoginResult, Course, Category, Tag, User, Order, Guestbook, Comment, NavItem, Banner, DashboardStats, ChartData, Ticket, TicketDetail, Settings, Admin } from '@/types'

// POST /api/v1/admin/login
export const login = (data: LoginParams) => request.post<unknown, ApiResponse<LoginResult>>('/admin/login', data)

// GET /api/v1/admin/dashboard/stats
export const getDashboardStats = () => request.get<unknown, ApiResponse<DashboardStats>>('/admin/dashboard/stats')
// GET /api/v1/admin/dashboard/charts?period=xxx
export const getDashboardCharts = (period: string) => request.get<unknown, ApiResponse<ChartData>>('/admin/dashboard/charts', { params: { period } })

// 课程管理
// GET /api/v1/admin/courses
export const getCourses = (params: Record<string, unknown>) => request.get<unknown, PageResponse<Course>>('/admin/courses', { params })
// POST /api/v1/admin/courses
export const createCourse = (data: Partial<Course>) => request.post<unknown, ApiResponse<Course>>('/admin/courses', data)
// PUT /api/v1/admin/courses/:id
export const updateCourse = (id: number, data: Partial<Course>) => request.put<unknown, ApiResponse<Course>>(`/admin/courses/${id}`, data)
// DELETE /api/v1/admin/courses/:id
export const deleteCourse = (id: number) => request.delete<unknown, ApiResponse<null>>(`/admin/courses/${id}`)
// PUT /api/v1/admin/courses/:id/status
export const updateCourseStatus = (id: number, status: string) => request.put<unknown, ApiResponse<null>>(`/admin/courses/${id}/status`, { status })
// POST /api/v1/admin/courses/batch
export const batchCourses = (ids: number[], action: string) => request.post<unknown, ApiResponse<null>>('/admin/courses/batch', { ids, action })

// 分类管理
// GET /api/v1/admin/categories
export const getCategories = (params?: Record<string, unknown>) => request.get<unknown, ApiResponse<Category[]>>('/admin/categories', { params })
// POST /api/v1/admin/categories
export const createCategory = (data: Partial<Category>) => request.post<unknown, ApiResponse<Category>>('/admin/categories', data)
// PUT /api/v1/admin/categories/:id
export const updateCategory = (id: number, data: Partial<Category>) => request.put<unknown, ApiResponse<Category>>(`/admin/categories/${id}`, data)
// DELETE /api/v1/admin/categories/:id
export const deleteCategory = (id: number) => request.delete<unknown, ApiResponse<null>>(`/admin/categories/${id}`)
// PUT /api/v1/admin/categories/:id/status
export const updateCategoryStatus = (id: number, status: string) => request.put<unknown, ApiResponse<null>>(`/admin/categories/${id}/status`, { status })

// 标签管理
// GET /api/v1/admin/tags
export const getTags = (params?: Record<string, unknown>) => request.get<unknown, ApiResponse<Tag[]>>('/admin/tags', { params })
// POST /api/v1/admin/tags
export const createTag = (data: { name: string }) => request.post<unknown, ApiResponse<Tag>>('/admin/tags', data)
// PUT /api/v1/admin/tags/:id
export const updateTag = (id: number, data: { name: string }) => request.put<unknown, ApiResponse<Tag>>(`/admin/tags/${id}`, data)
// DELETE /api/v1/admin/tags/:id
export const deleteTag = (id: number) => request.delete<unknown, ApiResponse<null>>(`/admin/tags/${id}`)

// 用户管理
// GET /api/v1/admin/users
export const getUsers = (params: Record<string, unknown>) => request.get<unknown, PageResponse<User>>('/admin/users', { params })
// GET /api/v1/admin/users/:id
export const getUser = (id: number) => request.get<unknown, ApiResponse<User>>(`/admin/users/${id}`)
// PUT /api/v1/admin/users/:id/status
export const updateUserStatus = (id: number, status: string) => request.put<unknown, ApiResponse<null>>(`/admin/users/${id}/status`, { status })
// POST /api/v1/admin/users/:id/recharge
export const rechargeUser = (id: number, amount: number) => request.post<unknown, ApiResponse<null>>(`/admin/users/${id}/recharge`, { amount })

// 订单管理
// GET /api/v1/admin/orders
export const getOrders = (params: Record<string, unknown>) => request.get<unknown, PageResponse<Order>>('/admin/orders', { params })
// GET /api/v1/admin/orders/:id
export const getOrder = (id: number) => request.get<unknown, ApiResponse<Order>>(`/admin/orders/${id}`)
// POST /api/v1/admin/orders/:id/refund
export const refundOrder = (id: number) => request.post<unknown, ApiResponse<null>>(`/admin/orders/${id}/refund`)

// 留言管理
// GET /api/v1/admin/guestbook
export const getGuestbook = (params: Record<string, unknown>) => request.get<unknown, PageResponse<Guestbook>>('/admin/guestbook', { params })
// PUT /api/v1/admin/guestbook/:id/status
export const updateGuestbookStatus = (id: number, status: string) => request.put<unknown, ApiResponse<null>>(`/admin/guestbook/${id}/status`, { status })
// PUT /api/v1/admin/guestbook/:id/pin
export const pinGuestbook = (id: number, pinned: boolean) => request.put<unknown, ApiResponse<null>>(`/admin/guestbook/${id}/pin`, { pinned })
// DELETE /api/v1/admin/guestbook/:id
export const deleteGuestbook = (id: number) => request.delete<unknown, ApiResponse<null>>(`/admin/guestbook/${id}`)

// 评论管理
// GET /api/v1/admin/comments
export const getComments = (params: Record<string, unknown>) => request.get<unknown, PageResponse<Comment>>('/admin/comments', { params })
// PUT /api/v1/admin/comments/:id/status
export const updateCommentStatus = (id: number, status: string) => request.put<unknown, ApiResponse<null>>(`/admin/comments/${id}/status`, { status })
// DELETE /api/v1/admin/comments/:id
export const deleteComment = (id: number) => request.delete<unknown, ApiResponse<null>>(`/admin/comments/${id}`)

// 首页配置 - 金刚区
// GET /api/v1/admin/nav-items
export const getNavItems = () => request.get<unknown, ApiResponse<NavItem[]>>('/admin/nav-items')
// POST /api/v1/admin/nav-items
export const createNavItem = (data: Partial<NavItem>) => request.post<unknown, ApiResponse<NavItem>>('/admin/nav-items', data)
// PUT /api/v1/admin/nav-items/:id
export const updateNavItem = (id: number, data: Partial<NavItem>) => request.put<unknown, ApiResponse<NavItem>>(`/admin/nav-items/${id}`, data)
// DELETE /api/v1/admin/nav-items/:id
export const deleteNavItem = (id: number) => request.delete<unknown, ApiResponse<null>>(`/admin/nav-items/${id}`)

// 首页配置 - Banner
// GET /api/v1/admin/banners
export const getBanners = () => request.get<unknown, ApiResponse<Banner[]>>('/admin/banners')
// POST /api/v1/admin/banners
export const createBanner = (data: Partial<Banner>) => request.post<unknown, ApiResponse<Banner>>('/admin/banners', data)
// PUT /api/v1/admin/banners/:id
export const updateBanner = (id: number, data: Partial<Banner>) => request.put<unknown, ApiResponse<Banner>>(`/admin/banners/${id}`, data)
// DELETE /api/v1/admin/banners/:id
export const deleteBanner = (id: number) => request.delete<unknown, ApiResponse<null>>(`/admin/banners/${id}`)

// POST /api/v1/upload/image
export const uploadImage = (file: File) => {
  const form = new FormData()
  form.append('file', file)
  return request.post<unknown, ApiResponse<{ url: string }>>('/upload/image', form)
}

// 工单管理
export const getTickets = (params: Record<string, unknown>) => request.get<unknown, PageResponse<Ticket>>('/admin/tickets', { params })
export const getTicket = (id: number) => request.get<unknown, ApiResponse<TicketDetail>>(`/admin/tickets/${id}`)
export const replyTicket = (id: number, content: string) => request.post<unknown, ApiResponse<null>>(`/admin/tickets/${id}/reply`, { content })
export const updateTicketStatus = (id: number, status: string) => request.put<unknown, ApiResponse<null>>(`/admin/tickets/${id}/status`, { status })

// 系统设置
export const getSettings = () => request.get<unknown, ApiResponse<Settings>>('/admin/settings')
export const updateSettings = (data: Partial<Settings>) => request.put<unknown, ApiResponse<null>>('/admin/settings', data)

// 管理员管理
export const getAdmins = (params?: Record<string, unknown>) => request.get<unknown, PageResponse<Admin>>('/admin/admins', { params })
export const createAdmin = (data: { username: string; password: string; role: string }) => request.post<unknown, ApiResponse<Admin>>('/admin/admins', data)
export const updateAdmin = (id: number, data: Partial<Admin & { password?: string }>) => request.put<unknown, ApiResponse<null>>(`/admin/admins/${id}`, data)
export const deleteAdmin = (id: number) => request.delete<unknown, ApiResponse<null>>(`/admin/admins/${id}`)
