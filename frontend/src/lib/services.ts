import api from '@/lib/api';
import type {
  AuthResponse, CaptchaResponse, CategoryBrief, CoinBalance,
  CommentItem, CourseDetail, CourseListItem, DownloadItem,
  FavoriteItem, GuestbookItem, NavItem, Banner, OrderItem,
  PaginatedList, TagBrief, Ticket, VipPackage, VipStatus, UserResponse,
} from '@/types';

// Auth
export const getCaptcha = () => api.post<CaptchaResponse>('/auth/captcha').then(r => r.data);
export const sendSMS = (data: { phone: string; captcha: string; captcha_key: string }) =>
  api.post('/auth/sms/send', data);
export const register = (data: { phone: string; sms_code: string; password: string }) =>
  api.post<AuthResponse>('/auth/register', data).then(r => r.data);
export const login = (data: { phone: string; password: string; captcha: string; captcha_key: string }) =>
  api.post<AuthResponse>('/auth/login', data).then(r => r.data);

// Courses
export const getLatestCourses = () =>
  api.get<CourseListItem[]>('/courses/latest').then(r => r.data);
export const getCourses = (params: { page?: number; pageSize?: number; categoryId?: number; sort?: string }) =>
  api.get<PaginatedList<CourseListItem>>('/courses', { params }).then(r => r.data);
export const getCourseDetail = (slug: string) =>
  api.get<CourseDetail>(`/courses/${slug}`).then(r => r.data);
export const searchCourses = (params: { q: string; page?: number; pageSize?: number }) =>
  api.get<PaginatedList<CourseListItem>>('/search', { params }).then(r => r.data);

// Categories & Tags
export const getCategories = () => api.get<CategoryBrief[]>('/categories').then(r => r.data);
export const getTags = () => api.get<TagBrief[]>('/tags').then(r => r.data);

// Home
export const getNavItems = () => api.get<NavItem[]>('/home/nav-items').then(r => r.data);
export const getBanners = () => api.get<Banner[]>('/home/banners').then(r => r.data);

// VIP
export const getVipPackages = () => api.get<VipPackage[]>('/vip/packages').then(r => r.data);
export const getVipStatus = () => api.get<VipStatus>('/vip/status').then(r => r.data);
export const purchaseVip = (packageId: number) => api.post('/vip/purchase', { package_id: packageId });

// Guestbook
export const getGuestbook = (params: { page?: number; pageSize?: number }) =>
  api.get<PaginatedList<GuestbookItem>>('/guestbook', { params }).then(r => r.data);
export const createGuestbook = (content: string) => api.post('/guestbook', { content });
export const likeGuestbook = (id: number) => api.post(`/guestbook/${id}/like`);

// Comments
export const getComments = (params: { target_type: string; target_id: number; page?: number }) =>
  api.get<PaginatedList<CommentItem>>('/comments', { params }).then(r => r.data);
export const createComment = (data: { target_type: string; target_id: number; content: string; parent_id?: number }) =>
  api.post('/comments', data);

// User
export const getProfile = () => api.get<UserResponse>('/user/profile').then(r => r.data);
export const changePassword = (data: { old_password: string; new_password: string }) =>
  api.put('/user/password', data);
export const getUserOrders = (params?: { page?: number }) =>
  api.get<PaginatedList<OrderItem>>('/user/orders', { params }).then(r => r.data);
export const getUserDownloads = (params?: { page?: number }) =>
  api.get<PaginatedList<DownloadItem>>('/user/downloads', { params }).then(r => r.data);
export const getUserFavorites = (params?: { page?: number }) =>
  api.get<PaginatedList<FavoriteItem>>('/user/favorites', { params }).then(r => r.data);
export const addFavorite = (courseId: number) => api.post('/user/favorites', { course_id: courseId });
export const removeFavorite = (courseId: number) => api.delete(`/user/favorites/${courseId}`);

// Coins
export const getCoinBalance = () => api.get<CoinBalance>('/coins/balance').then(r => r.data);
export const recharge = (data: { amount: number; pay_method: string }) => api.post('/coins/recharge', data);

// Course actions
export const likeCourse = (id: number) => api.post(`/courses/${id}/like`);
export const downloadCourse = (id: number) => api.post(`/courses/${id}/download`);
export const purchaseCourse = (courseId: number) =>
  api.post('/orders', { type: 'course', target_id: courseId });

// Orders
export const createOrder = (data: { type: string; target_id?: number; amount?: number }) =>
  api.post('/orders', data);

// Forgot Password
export const forgotPassword = (data: { phone: string; sms_code: string; new_password: string }) =>
  api.post('/auth/forgot-password', data);

// Tickets
export const getTickets = (params?: { page?: number }) =>
  api.get<PaginatedList<Ticket>>('/tickets', { params }).then(r => r.data);
export const createTicket = (data: { title: string; content: string }) =>
  api.post('/tickets', data);
export const getTicketDetail = (id: number) =>
  api.get<Ticket>(`/tickets/${id}`).then(r => r.data);
export const replyTicket = (id: number, content: string) =>
  api.post(`/tickets/${id}/reply`, { content });
