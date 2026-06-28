import request from '@/utils/request'
import type { ApiResponse, PageData, PageResponse, LoginParams, LoginResult, Course, Category, Tag, User, Order, Guestbook, Comment, NavItem, Banner, VipPackageConfig, DashboardStats, ChartData, Ticket, TicketDetail, Settings, Admin } from '@/types'

type AnyRecord = Record<string, any>

const apiAssetOrigin = (() => {
  const configured = import.meta.env.VITE_API_ORIGIN
  if (configured) return configured.replace(/\/$/, '')
  if (typeof window !== 'undefined' && ['localhost', '127.0.0.1'].includes(window.location.hostname)) {
    return 'http://localhost:8080'
  }
  return 'https://api.zioran.com'
})()

export const assetUrl = (url?: string) => {
  if (!url) return ''
  if (/^(https?:)?\/\//.test(url) || url.startsWith('data:') || url.startsWith('blob:')) return url
  if (url.startsWith('/uploads/')) return `${apiAssetOrigin}${url}`
  return url
}

const slugify = (value?: string, fallback = 'item') => {
  const slug = (value || '')
    .trim()
    .toLowerCase()
    .replace(/\s+/g, '-')
    .replace(/[^a-z0-9-]/g, '')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')
  return slug || `${fallback}-${Date.now()}`
}

const mapApi = async <T, U>(promise: Promise<ApiResponse<T>>, mapper: (data: T) => U): Promise<ApiResponse<U>> => {
  const res = await promise
  return { ...res, data: mapper(res.data) }
}

const mapPage = <T extends AnyRecord, U>(data: AnyRecord, mapper: (item: T) => U): PageData<U> => ({
  items: Array.isArray(data.items) ? data.items.map(mapper) : [],
  total: Number(data.total || 0),
  page: Number(data.page || 1),
  pageSize: Number(data.pageSize || 20),
  totalPages: Number(data.totalPages || 0),
})

const adaptCourse = (item: AnyRecord): Course => ({
  ...item,
  coverImage: assetUrl(item.coverImage || item.cover || ''),
  categoryId: item.categoryId ?? item.category?.id ?? 0,
  categoryName: item.categoryName || item.category?.name || '-',
  vipPrice: item.vipPrice ?? 0,
  qualityLabel: item.qualityLabel || '',
  detailTitle: item.detailTitle || '',
  detailSubtitle: item.detailSubtitle || '',
  detailImages: Array.isArray(item.detailImages) ? item.detailImages.map(assetUrl) : [],
  tags: Array.isArray(item.tags) ? item.tags : [],
  resources: Array.isArray(item.resources)
    ? item.resources.map((r: AnyRecord) => ({ link: r.link || r.url || '', code: r.code || r.password || '' }))
    : [],
  createdAt: item.createdAt || '',
  updatedAt: item.updatedAt || '',
} as Course)

const adaptCategory = (item: AnyRecord): Category => ({
  ...item,
  parentId: item.parentId ?? 0,
  sort: item.sort ?? item.sortOrder ?? 0,
  status: item.status || (item.isActive === false ? 'inactive' : 'active'),
  createdAt: item.createdAt || '',
} as Category)

const adaptTag = (item: AnyRecord): Tag => ({ ...item, createdAt: item.createdAt || '' } as Tag)

const adaptUser = (item: AnyRecord): User => ({
  ...item,
  nickname: item.nickname || item.username || '-',
  avatar: item.avatar || item.avatarUrl || '',
  isVip: item.isVip ?? item.is_vip ?? false,
  balance: item.balance ?? 0,
  vipExpireAt: item.vipExpireAt || item.vip_expire_at || '',
  purchasedCount: item.purchasedCount ?? item.purchased_count ?? 0,
  favoriteCount: item.favoriteCount ?? item.favorite_count ?? 0,
  createdAt: item.createdAt || item.created_at || '',
} as User)

const adaptOrder = (item: AnyRecord): Order => ({
  ...item,
  productName: item.productName || item.targetName || '-',
  userName: item.userName || '-',
  payMethod: item.payMethod || '-',
  createdAt: item.createdAt || '',
  paidAt: item.paidAt || '',
} as Order)

const adaptGuestbook = (item: AnyRecord): Guestbook => ({
  ...item,
  userName: item.userName || item.username || item.user?.username || '-',
  userAvatar: item.userAvatar || item.avatar || item.user?.avatarUrl || '',
  likes: item.likes ?? item.likeCount ?? 0,
  pinned: item.pinned ?? item.isPinned ?? false,
  createdAt: item.createdAt || '',
} as Guestbook)

const adaptComment = (item: AnyRecord): Comment => ({
  ...item,
  userName: item.userName || item.username || item.user?.username || '-',
  targetName: item.targetName || `${item.targetType || '内容'} #${item.targetId || '-'}`,
  createdAt: item.createdAt || '',
} as Comment)

const adaptNavItem = (item: AnyRecord): NavItem => ({
  ...item,
  subtitle: item.subtitle || '',
  icon: assetUrl(item.icon || ''),
  link: item.link || item.url || '',
  categoryId: item.categoryId ?? item.category_id ?? null,
  sort: item.sort ?? item.sortOrder ?? 0,
  status: item.status || (item.isActive === false ? 'inactive' : 'active'),
} as NavItem)

const adaptBanner = (item: AnyRecord): Banner => ({
  ...item,
  image: assetUrl(item.image || item.imageUrl || ''),
  link: item.link || item.linkUrl || '',
  placement: item.placement || 'home',
  backgroundColor: item.backgroundColor || '',
  sort: item.sort ?? item.sortOrder ?? 0,
  status: item.status || (item.isActive === false ? 'inactive' : 'active'),
} as Banner)

const adaptVipPackage = (item: AnyRecord): VipPackageConfig => ({
  ...item,
  originalPrice: item.originalPrice ?? item.original_price ?? 0,
  benefits: item.benefits || '',
  isActive: item.isActive ?? item.is_active ?? true,
  sortOrder: item.sortOrder ?? item.sort_order ?? 0,
} as VipPackageConfig)

const adaptTicket = (item: AnyRecord): Ticket => ({
  ...item,
  userName: item.userName || item.username || '-',
  subject: item.subject || item.title || '',
  attachments: Array.isArray(item.attachments) ? item.attachments.map(assetUrl) : [],
  createdAt: item.createdAt || '',
  updatedAt: item.updatedAt || '',
} as Ticket)

const adaptTicketDetail = (item: AnyRecord): TicketDetail => ({
  ...adaptTicket(item),
  content: item.content || '',
  attachments: Array.isArray(item.attachments) ? item.attachments.map(assetUrl) : [],
  replies: Array.isArray(item.replies)
    ? item.replies.map((r: AnyRecord) => ({
      ...r,
      ticketId: r.ticketId ?? item.id,
      userName: r.userName || r.username || '-',
      createdAt: r.createdAt || '',
    }))
    : [],
} as TicketDetail)

const stripAssetOrigin = (url: string) => {
  if (!url) return ''
  try {
    const u = new URL(url)
    if (u.pathname.startsWith('/uploads/')) return u.pathname
  } catch { /* not a full URL */ }
  return url
}

const toCoursePayload = (data: AnyRecord) => ({
  title: data.title,
  subtitle: data.subtitle || '',
  slug: data.slug || slugify(data.title, 'course'),
  category_id: data.categoryId,
  quality_label: data.qualityLabel || '',
  cover_image: stripAssetOrigin(data.coverImage || ''),
  content: data.content || '',
  detail_title: data.detailTitle || '',
  detail_subtitle: data.detailSubtitle || '',
  price: data.price ?? 0,
  vip_price: data.vipPrice ?? 0,
  tag_ids: Array.isArray(data.tags) ? data.tags.map((t: AnyRecord | number) => typeof t === 'number' ? t : t.id).filter(Boolean) : [],
  resources: Array.isArray(data.resources)
    ? data.resources.map((r: AnyRecord, index: number) => ({
      name: r.name || `资源${index + 1}`,
      url: r.url || r.link || '',
      password: r.password || r.code || '',
      sort_order: index,
    })).filter((r: AnyRecord) => r.url)
    : [],
})

const toCategoryPayload = (data: AnyRecord) => ({
  name: data.name,
  slug: data.slug || slugify(data.name, 'category'),
  parent_id: data.parentId ? data.parentId : null,
  sort_order: data.sort ?? data.sortOrder ?? 0,
  is_active: data.status ? data.status === 'active' : data.isActive ?? true,
})

const toTagPayload = (data: AnyRecord) => ({ name: data.name, slug: data.slug || slugify(data.name, 'tag') })

const toNavItemPayload = (data: AnyRecord) => ({
  title: data.title,
  subtitle: data.subtitle || '',
  icon: data.icon || '',
  url: data.url || data.link || '',
  category_id: data.categoryId || null,
  sort_order: data.sort ?? data.sortOrder ?? 0,
  is_active: data.status ? data.status === 'active' : data.isActive ?? true,
})

const toBannerPayload = (data: AnyRecord) => ({
  title: data.title || '',
  image_url: stripAssetOrigin(data.imageUrl || data.image || ''),
  link_url: data.linkUrl || data.link || '',
  placement: data.placement || 'home',
  background_color: data.backgroundColor || data.background_color || '',
  sort_order: data.sort ?? data.sortOrder ?? 0,
  is_active: data.status ? data.status === 'active' : data.isActive ?? true,
})

const toVipPackagePayload = (data: Partial<VipPackageConfig>) => ({
  name: data.name,
  price: data.price ?? 0,
  original_price: data.originalPrice ?? 0,
  benefits: data.benefits || '',
  is_active: data.isActive ?? true,
  sort_order: data.sortOrder ?? 0,
})

const adaptSettings = (data: AnyRecord): Settings => ({
  siteName: data.siteName || '',
  siteDescription: data.siteDescription || '',
  contactPhone: data.contactPhone || '',
  contactEmail: data.contactEmail || '',
  vipMonthlyPrice: Number(data.vipMonthlyPrice || 0),
  vipYearlyPrice: Number(data.vipYearlyPrice || 0),
  withdrawMinAmount: Number(data.withdrawMinAmount || 0),
  commissionRate: Number(data.commissionRate || 0),
  coinRechargeRatio: Number(data.coinRechargeRatio || 1),
  coinRechargeAmounts: data.coinRechargeAmounts || '10,50,100,200,500,1000',
})

const toSettingsPayload = (data: Partial<Settings>) => Object.fromEntries(
  Object.entries(data).map(([key, value]) => [key, value == null ? '' : String(value)])
)

// POST /api/v1/admin/login
export const login = (data: LoginParams) => request.post<unknown, ApiResponse<LoginResult>>('/admin/login', data)

// GET /api/v1/admin/dashboard/stats
export const getDashboardStats = () => request.get<unknown, ApiResponse<DashboardStats>>('/admin/dashboard/stats')
// GET /api/v1/admin/dashboard/charts?period=xxx
export const getDashboardCharts = (period: string) => request.get<unknown, ApiResponse<ChartData>>('/admin/dashboard/charts', { params: { period } })

// 课程管理
// GET /api/v1/admin/courses
export const getCourses = (params: Record<string, unknown>) =>
  mapApi(request.get<unknown, PageResponse<Course>>('/admin/courses', { params }), data => mapPage(data, adaptCourse))
export const getCourse = async (id: number): Promise<ApiResponse<Course | null>> => {
  let page = 1
  const pageSize = 100
  while (true) {
    const res = await getCourses({ page, pageSize })
    const course = res.data.items.find(item => item.id === id)
    if (course || page >= res.data.totalPages || res.data.items.length === 0) {
      return { ...res, data: course || null }
    }
    page += 1
  }
}
// POST /api/v1/admin/courses
export const createCourse = async (data: Partial<Course>) => {
  const res = await mapApi(request.post<unknown, ApiResponse<Course>>('/admin/courses', toCoursePayload(data)), adaptCourse)
  if (data.status && data.status !== 'draft' && res.data.id) {
    await request.put<unknown, ApiResponse<null>>(`/admin/courses/${res.data.id}/status`, { status: data.status })
  }
  return res
}
// PUT /api/v1/admin/courses/:id
export const updateCourse = async (id: number, data: Partial<Course>) => {
  const res = await mapApi(request.put<unknown, ApiResponse<Course>>(`/admin/courses/${id}`, toCoursePayload(data)), adaptCourse)
  if (data.status) await request.put<unknown, ApiResponse<null>>(`/admin/courses/${id}/status`, { status: data.status })
  return res
}
// DELETE /api/v1/admin/courses/:id
export const deleteCourse = (id: number) => request.delete<unknown, ApiResponse<null>>(`/admin/courses/${id}`)
// PUT /api/v1/admin/courses/:id/status
export const updateCourseStatus = (id: number, status: string) => request.put<unknown, ApiResponse<null>>(`/admin/courses/${id}/status`, { status })
// POST /api/v1/admin/courses/batch
export const batchCourses = (ids: number[], action: string) => request.post<unknown, ApiResponse<null>>('/admin/courses/batch', { ids, action })

// 分类管理
// GET /api/v1/admin/categories
export const getCategories = (params?: Record<string, unknown>) =>
  mapApi(request.get<unknown, ApiResponse<Category[]>>('/admin/categories', { params }), data => data.map(adaptCategory))
// POST /api/v1/admin/categories
export const createCategory = (data: Partial<Category>) =>
  mapApi(request.post<unknown, ApiResponse<Category>>('/admin/categories', toCategoryPayload(data)), adaptCategory)
// PUT /api/v1/admin/categories/:id
export const updateCategory = (id: number, data: Partial<Category>) =>
  mapApi(request.put<unknown, ApiResponse<Category>>(`/admin/categories/${id}`, toCategoryPayload(data)), adaptCategory)
// DELETE /api/v1/admin/categories/:id
export const deleteCategory = (id: number) => request.delete<unknown, ApiResponse<null>>(`/admin/categories/${id}`)
// PUT /api/v1/admin/categories/:id/status
export const updateCategoryStatus = (id: number, status: string) => request.put<unknown, ApiResponse<null>>(`/admin/categories/${id}/status`, { is_active: status === 'active' })

// 标签管理
// GET /api/v1/admin/tags
export const getTags = (params?: Record<string, unknown>) =>
  mapApi(request.get<unknown, ApiResponse<Tag[]>>('/admin/tags', { params }), data => data.map(adaptTag))
// POST /api/v1/admin/tags
export const createTag = (data: { name: string }) =>
  mapApi(request.post<unknown, ApiResponse<Tag>>('/admin/tags', toTagPayload(data)), adaptTag)
// PUT /api/v1/admin/tags/:id
export const updateTag = (id: number, data: { name: string }) =>
  mapApi(request.put<unknown, ApiResponse<Tag>>(`/admin/tags/${id}`, toTagPayload(data)), adaptTag)
// DELETE /api/v1/admin/tags/:id
export const deleteTag = (id: number) => request.delete<unknown, ApiResponse<null>>(`/admin/tags/${id}`)

// 用户管理
// GET /api/v1/admin/users
export const getUsers = (params: Record<string, unknown>) =>
  mapApi(request.get<unknown, PageResponse<User>>('/admin/users', { params }), data => mapPage(data, adaptUser))
// GET /api/v1/admin/users/:id
export const getUser = (id: number) => mapApi(request.get<unknown, ApiResponse<User>>(`/admin/users/${id}`), adaptUser)
// PUT /api/v1/admin/users/:id/status
export const updateUserStatus = (id: number, status: string) => request.put<unknown, ApiResponse<null>>(`/admin/users/${id}/status`, { status })
// POST /api/v1/admin/users/:id/recharge
export const rechargeUser = (id: number, amount: number) => request.post<unknown, ApiResponse<null>>(`/admin/users/${id}/recharge`, { amount })

// 订单管理
// GET /api/v1/admin/orders
export const getOrders = (params: Record<string, unknown>) =>
  mapApi(request.get<unknown, PageResponse<Order>>('/admin/orders', { params }), data => mapPage(data, adaptOrder))
// GET /api/v1/admin/orders/:id
export const getOrder = (id: number) => mapApi(request.get<unknown, ApiResponse<Order>>(`/admin/orders/${id}`), adaptOrder)
// POST /api/v1/admin/orders/:id/refund
export const refundOrder = (id: number) => request.post<unknown, ApiResponse<null>>(`/admin/orders/${id}/refund`)

// 留言管理
// GET /api/v1/admin/guestbook
export const getGuestbook = (params: Record<string, unknown>) =>
  mapApi(request.get<unknown, PageResponse<Guestbook>>('/admin/guestbook', { params }), data => mapPage(data, adaptGuestbook))
// PUT /api/v1/admin/guestbook/:id/status
export const updateGuestbookStatus = (id: number, status: string) => request.put<unknown, ApiResponse<null>>(`/admin/guestbook/${id}/status`, { status })
// PUT /api/v1/admin/guestbook/:id/pin
export const pinGuestbook = (id: number, pinned: boolean) => request.put<unknown, ApiResponse<null>>(`/admin/guestbook/${id}/pin`, { pinned })
// DELETE /api/v1/admin/guestbook/:id
export const deleteGuestbook = (id: number) => request.delete<unknown, ApiResponse<null>>(`/admin/guestbook/${id}`)

// 评论管理
// GET /api/v1/admin/comments
export const getComments = (params: Record<string, unknown>) =>
  mapApi(request.get<unknown, PageResponse<Comment>>('/admin/comments', { params }), data => mapPage(data, adaptComment))
// PUT /api/v1/admin/comments/:id/status
export const updateCommentStatus = (id: number, status: string) => request.put<unknown, ApiResponse<null>>(`/admin/comments/${id}/status`, { status })
// DELETE /api/v1/admin/comments/:id
export const deleteComment = (id: number) => request.delete<unknown, ApiResponse<null>>(`/admin/comments/${id}`)

// 首页配置 - 金刚区
// GET /api/v1/admin/nav-items
export const getNavItems = () => mapApi(request.get<unknown, ApiResponse<NavItem[]>>('/admin/nav-items'), data => data.map(adaptNavItem))
// POST /api/v1/admin/nav-items
export const createNavItem = (data: Partial<NavItem>) =>
  mapApi(request.post<unknown, ApiResponse<NavItem>>('/admin/nav-items', toNavItemPayload(data)), adaptNavItem)
// PUT /api/v1/admin/nav-items/:id
export const updateNavItem = (id: number, data: Partial<NavItem>) =>
  mapApi(request.put<unknown, ApiResponse<NavItem>>(`/admin/nav-items/${id}`, toNavItemPayload(data)), adaptNavItem)
// DELETE /api/v1/admin/nav-items/:id
export const deleteNavItem = (id: number) => request.delete<unknown, ApiResponse<null>>(`/admin/nav-items/${id}`)

// 首页配置 - Banner
// GET /api/v1/admin/banners
export const getBanners = () => mapApi(request.get<unknown, ApiResponse<Banner[]>>('/admin/banners'), data => data.map(adaptBanner))
// POST /api/v1/admin/banners
export const createBanner = (data: Partial<Banner>) =>
  mapApi(request.post<unknown, ApiResponse<Banner>>('/admin/banners', toBannerPayload(data)), adaptBanner)
// PUT /api/v1/admin/banners/:id
export const updateBanner = (id: number, data: Partial<Banner>) =>
  mapApi(request.put<unknown, ApiResponse<Banner>>(`/admin/banners/${id}`, toBannerPayload(data)), adaptBanner)
// DELETE /api/v1/admin/banners/:id
export const deleteBanner = (id: number) => request.delete<unknown, ApiResponse<null>>(`/admin/banners/${id}`)

// VIP 套餐配置
export const getVipPackagesAdmin = () =>
  mapApi(request.get<unknown, ApiResponse<VipPackageConfig[]>>('/admin/vip/packages'), data => data.map(adaptVipPackage))
export const updateVipPackage = (id: number, data: Partial<VipPackageConfig>) =>
  mapApi(request.put<unknown, ApiResponse<VipPackageConfig>>(`/admin/vip/packages/${id}`, toVipPackagePayload(data)), adaptVipPackage)

// POST /api/v1/upload/image
export const uploadImage = (file: File) => {
  const form = new FormData()
  form.append('file', file)
  return mapApi(request.post<unknown, ApiResponse<{ url: string }>>('/upload/image', form), data => ({ url: assetUrl(data.url) }))
}

// 工单管理
export const getTickets = (params: Record<string, unknown>) =>
  mapApi(request.get<unknown, PageResponse<Ticket>>('/admin/tickets', { params }), data => mapPage(data, adaptTicket))
export const getTicket = (id: number) => mapApi(request.get<unknown, ApiResponse<TicketDetail>>(`/admin/tickets/${id}`), adaptTicketDetail)
export const replyTicket = (id: number, content: string) => request.post<unknown, ApiResponse<null>>(`/admin/tickets/${id}/reply`, { content })
export const updateTicketStatus = (id: number, status: string) => request.put<unknown, ApiResponse<null>>(`/admin/tickets/${id}/status`, { status })

// 系统设置
export const getSettings = () => mapApi(request.get<unknown, ApiResponse<Settings>>('/admin/settings'), adaptSettings)
export const updateSettings = (data: Partial<Settings>) => request.put<unknown, ApiResponse<null>>('/admin/settings', toSettingsPayload(data))

// 管理员管理
export const getAdmins = async (params?: Record<string, unknown>): Promise<PageResponse<Admin>> => {
  const res = await request.get<unknown, ApiResponse<Admin[] | AnyRecord>>('/admin/admins', { params })
  const list = Array.isArray(res.data) ? res.data : (res.data.items || [])
  const page = Number(params?.page || 1)
  const pageSize = Number(params?.pageSize || 20)
  return {
    ...res,
    data: {
      items: list,
      total: Array.isArray(res.data) ? list.length : res.data.total,
      page,
      pageSize,
      totalPages: Math.ceil((Array.isArray(res.data) ? list.length : res.data.total || 0) / pageSize),
    },
  }
}
export const createAdmin = (data: { username: string; password: string; role: string }) => request.post<unknown, ApiResponse<Admin>>('/admin/admins', data)
export const updateAdmin = (id: number, data: Partial<Admin & { password?: string }>) => request.put<unknown, ApiResponse<null>>(`/admin/admins/${id}`, data)
export const deleteAdmin = (id: number) => request.delete<unknown, ApiResponse<null>>(`/admin/admins/${id}`)
