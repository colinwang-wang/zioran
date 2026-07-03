// API 通用响应
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface PageData<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

export type PageResponse<T> = ApiResponse<PageData<T>>

// 课程
export interface Course {
  id: number
  title: string
  subtitle: string
  slug: string
  coverImage: string
  price: number
  vipPrice: number
  categoryId: number
  categoryName: string
  tags: Tag[]
  status: 'published' | 'draft' | 'offline'
  qualityLabel: string
  detailTitle: string
  detailSubtitle: string
  detailImages: string[]
  resources: Resource[]
  createdAt: string
  updatedAt: string
}

export interface Resource {
  link: string
  code: string
}

// 分类
export interface Category {
  id: number
  name: string
  parentId: number
  sort: number
  status: 'active' | 'inactive'
  children?: Category[]
  createdAt: string
}

// 标签
export interface Tag {
  id: number
  name: string
  createdAt: string
}

// 用户
export interface User {
  id: number
  phone: string
  username: string
  nickname: string
  avatar: string
  balance: number
  isVip: boolean
  vipExpireAt: string
  status: 'active' | 'disabled'
  purchasedCount: number
  favoriteCount: number
  createdAt: string
}

// 订单
export interface Order {
  id: number
  orderNo: string
  userId: number
  userName: string
  targetName: string
  productName: string
  amount: number
  payMethod: string
  type: 'coin' | 'vip' | 'course' | 'coin_recharge' | 'vip_purchase' | 'course_purchase'
  status: 'pending' | 'paid' | 'refunded' | 'cancelled'
  createdAt: string
  paidAt: string
}

// 留言
export interface Guestbook {
  id: number
  userId: number
  userName: string
  userAvatar: string
  content: string
  likes: number
  pinned: boolean
  status: 'visible' | 'hidden'
  createdAt: string
}

// 评论
export interface Comment {
  id: number
  userId: number
  userName: string
  content: string
  targetType: string
  targetId: number
  targetName: string
  status: 'visible' | 'hidden'
  createdAt: string
}

// 金刚区导航
export interface NavItem {
  id: number
  title: string
  subtitle: string
  icon: string
  link: string
  categoryId?: number | null
  sort: number
  status?: 'active' | 'inactive'
  createdAt?: string
}

// Banner
export interface Banner {
  id: number
  title: string
  image: string
  link: string
  placement: 'home' | 'vip'
  backgroundColor: string
  sort: number
  status: 'active' | 'inactive'
}

export interface VipPackageConfig {
  id: number
  name: string
  price: number
  originalPrice: number
  benefits: string
  isActive: boolean
  sortOrder: number
}

// 看板统计
export interface DashboardStats {
  totalUsers: number
  totalCourses: number
  totalOrders: number
  todayRevenue: number
  userGrowth: number
  courseGrowth: number
  orderGrowth: number
  revenueGrowth: number
  totalFavorites?: number
  pendingOrders?: number
  recentNewUsers?: number
}

export interface ChartData {
  labels: string[]
  datasets: { label: string; data: number[] }[]
}

// 工单
export interface Ticket {
  id: number
  userId: number
  userName: string
  subject: string
  content: string
  status: 'pending' | 'processing' | 'replied' | 'closed'
  attachments?: string[]
  createdAt: string
  updatedAt: string
}

export interface TicketReply {
  id: number
  ticketId: number
  userId: number
  userName: string
  content: string
  isAdmin: boolean
  createdAt: string
}

export interface TicketDetail extends Ticket {
  replies: TicketReply[]
}

// 系统设置
export interface Settings {
  siteName: string
  siteDescription: string
  contactPhone: string
  contactEmail: string
  vipMonthlyPrice: number
  vipYearlyPrice: number
  withdrawMinAmount: number
  commissionRate: number
  coinRechargeRatio: number
  coinRechargeAmounts: string
}

// 管理员
export interface Admin {
  id: number
  username: string
  role: string
  status: 'active' | 'disabled'
  createdAt: string
  lastLoginAt: string
}

// 管理员登录
export interface LoginParams {
  username: string
  password: string
}

export interface LoginResult {
  token: string
  admin: { id: number; username: string; role: string }
}
