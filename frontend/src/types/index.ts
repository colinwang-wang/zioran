export interface PaginatedList<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export interface CategoryBrief {
  id: number;
  name: string;
  slug: string;
}

export interface TagBrief {
  id: number;
  name: string;
  slug: string;
}

export interface CourseListItem {
  id: number;
  title: string;
  subtitle: string;
  slug: string;
  cover: string;
  category: CategoryBrief | null;
  price: number;
  vip_price: number;
  relative_time: string;
  published_at: string | null;
}

export interface ResourceItem {
  id: number;
  name: string;
  url: string;
  password: string;
}

export interface CourseBrief {
  slug: string;
  title: string;
}

export interface UserAccess {
  can_download: boolean;
  has_purchased: boolean;
  is_vip: boolean;
  is_favorited: boolean;
}

export interface CourseDetail {
  id: number;
  title: string;
  subtitle: string;
  slug: string;
  cover: string;
  content: string;
  detail_title: string;
  detail_subtitle: string;
  price: number;
  vip_price: number;
  quality_label: string;
  category: CategoryBrief | null;
  tags: TagBrief[];
  like_count: number;
  comment_count: number;
  published_at: string | null;
  prev_course: CourseBrief | null;
  next_course: CourseBrief | null;
  related_courses: CourseListItem[];
  user_access: UserAccess | null;
  resources: ResourceItem[];
}

export interface UserResponse {
  id: number;
  username: string;
  email: string;
  avatar: string;
  is_vip: boolean;
}

export interface AuthResponse {
  token: string;
  user: UserResponse;
}

export interface CaptchaResponse {
  captcha_key: string;
  captcha_image: string;
}

export interface NavItem {
  id: number;
  title: string;
  subtitle: string;
  icon: string;
  url: string;
  category_id?: number | null;
  sort_order: number;
}

export interface Banner {
  id: number;
  title: string;
  image_url: string;
  link_url: string;
  placement: string;
  background_color: string;
}

export interface GuestbookItem {
  id: number;
  user_id: number;
  username: string;
  avatar: string;
  content: string;
  like_count: number;
  is_pinned: boolean;
  is_liked: boolean;
  created_at: string;
}

export interface CommentItem {
  id: number;
  user_id: number;
  username: string;
  avatar: string;
  target_type: string;
  target_id: number;
  content: string;
  parent_id: number | null;
  status: string;
  created_at: string;
  children?: CommentItem[];
}

export interface CoinBalance {
  balance: number;
  total_earned: number;
  total_spent: number;
}

export interface CoinTransaction {
  id: number;
  user_id: number;
  type: string;
  amount: number;
  balance_after: number;
  description: string;
  order_id: number | null;
  created_at: string;
}

export interface RechargeResponse {
  order_id: number;
  order_no: string;
  pay_url: string;
  amount: number;
  coins: number;
}

export interface RechargeConfig {
  ratio: number;
  amounts: number[];
}

export interface CourseDownloadResponse {
  resources: ResourceItem[];
}

export interface VipStatus {
  is_vip: boolean;
  expires_at: string | null;
  package: string;
}

export interface OrderItem {
  id: number;
  order_no: string;
  type: string;
  target_name: string;
  amount: number;
  status: string;
  created_at: string;
  paid_at: string | null;
}

export interface DownloadItem {
  id: number;
  course_id: number;
  title: string;
  cover: string;
  order_no: string;
  amount: number;
  created_at: string;
}

export interface FavoriteItem {
  course_id: number;
  title: string;
  cover: string;
  slug: string;
  created_at: string;
}

export interface VipPackage {
  id: number;
  name: string;
  price: number;
  original_price: number;
  duration: string;
  features: string[];
}

export interface TicketReply {
  id: number;
  ticket_id: number;
  user_id: number;
  username: string;
  avatar: string;
  content: string;
  is_admin: boolean;
  created_at: string;
}

export interface Ticket {
  id: number;
  user_id: number;
  title: string;
  content: string;
  status: string;
  created_at: string;
  updated_at: string;
  attachments?: string[];
  replies?: TicketReply[];
}
