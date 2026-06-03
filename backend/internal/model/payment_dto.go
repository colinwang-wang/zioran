package model

import "time"

// Coin DTOs

type CoinBalanceResponse struct {
	Balance    int `json:"balance"`
	TotalEarned int `json:"total_earned"`
	TotalSpent int  `json:"total_spent"`
}

type RechargeRequest struct {
	Amount    int    `json:"amount" binding:"required"`
	PayMethod string `json:"pay_method" binding:"required"`
}

// VIP DTOs

type VipStatusResponse struct {
	IsVip     bool       `json:"is_vip"`
	ExpiresAt *time.Time `json:"expires_at"`
	Package   string     `json:"package"`
}

type VipPurchaseRequest struct {
	PackageID int `json:"package_id" binding:"required"`
}

// Order DTOs

type CreateOrderRequest struct {
	Type     string `json:"type" binding:"required"`     // course, vip, coin
	TargetID int    `json:"target_id"`
	Amount   int    `json:"amount"`
}

type OrderResponse struct {
	ID         int64      `json:"id"`
	OrderNo    string     `json:"order_no"`
	Type       string     `json:"type"`
	TargetName string     `json:"target_name"`
	Amount     int        `json:"amount"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	PaidAt     *time.Time `json:"paid_at"`
}

// Guestbook DTOs

type GuestbookCreateRequest struct {
	Content string `json:"content" binding:"required"`
}

type GuestbookResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	Content   string    `json:"content"`
	LikeCount int       `json:"like_count"`
	IsPinned  bool      `json:"is_pinned"`
	IsLiked   bool      `json:"is_liked"`
	CreatedAt time.Time `json:"created_at"`
}

// Comment DTOs

type CommentCreateRequest struct {
	TargetType string `json:"target_type" binding:"required"`
	TargetID   int64  `json:"target_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
	ParentID   *int64 `json:"parent_id"`
}

type CommentResponse struct {
	ID        int64             `json:"id"`
	UserID    int64             `json:"user_id"`
	Username  string            `json:"username"`
	Avatar    string            `json:"avatar"`
	Content   string            `json:"content"`
	ParentID  *int64            `json:"parent_id"`
	CreatedAt time.Time         `json:"created_at"`
	Children  []CommentResponse `json:"children,omitempty"`
}

// Home Config DTOs

type NavItemRequest struct {
	Title     string `json:"title" binding:"required"`
	Icon      string `json:"icon"`
	URL       string `json:"url" binding:"required"`
	SortOrder int    `json:"sort_order"`
	IsActive  *bool  `json:"is_active"`
}

type BannerRequest struct {
	Title     string `json:"title"`
	ImageURL  string `json:"image_url" binding:"required"`
	LinkURL   string `json:"link_url"`
	SortOrder int    `json:"sort_order"`
	IsActive  *bool  `json:"is_active"`
}

// User center DTOs

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type FavoriteRequest struct {
	CourseID int64 `json:"course_id" binding:"required"`
}

type DownloadResponse struct {
	ID        int64     `json:"id"`
	CourseID  int64     `json:"course_id"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	CreatedAt time.Time `json:"created_at"`
}

type FavoriteResponse struct {
	CourseID    int64      `json:"course_id"`
	Title       string     `json:"title"`
	Cover       string     `json:"cover"`
	Slug        string     `json:"slug"`
	CreatedAt   time.Time  `json:"created_at"`
}

// Admin DTOs

type AdminUserStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type AdminRechargeRequest struct {
	Amount      int    `json:"amount" binding:"required"`
	Description string `json:"description"`
}

type AdminGuestbookStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type AdminCommentStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type DashboardStats struct {
	TotalUsers    int64 `json:"total_users"`
	TotalCourses  int64 `json:"total_courses"`
	TotalOrders   int64 `json:"total_orders"`
	TotalRevenue  int64 `json:"total_revenue"`
	TodayUsers    int64 `json:"today_users"`
	TodayOrders   int64 `json:"today_orders"`
}

type CourseDownloadResponse struct {
	Resources []ResourceItem `json:"resources"`
}
