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
  nickname: string
  avatar: string
  balance: number
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
  productName: string
  amount: number
  payMethod: string
  type: 'coin_recharge' | 'vip_purchase' | 'course_purchase'
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
  icon: string
  link: string
  sort: number
}

// Banner
export interface Banner {
  id: number
  title: string
  image: string
  link: string
  sort: number
  status: 'active' | 'inactive'
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
}

export interface ChartData {
  labels: string[]
  datasets: { label: string; data: number[] }[]
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
