package model

import "time"

// Coin

type CoinAccount struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      int64     `json:"user_id" gorm:"not null;uniqueIndex"`
	Balance     int       `json:"balance" gorm:"not null;default:0"`
	TotalEarned int       `json:"total_earned" gorm:"not null;default:0"`
	TotalSpent  int       `json:"total_spent" gorm:"not null;default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (CoinAccount) TableName() string { return "coin_accounts" }

type CoinTransaction struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       int64     `json:"user_id" gorm:"not null"`
	Type         string    `json:"type" gorm:"size:20;not null"` // recharge, purchase, vip
	Amount       int       `json:"amount" gorm:"not null"`
	BalanceAfter int       `json:"balance_after" gorm:"not null"`
	Description  string    `json:"description" gorm:"size:200"`
	OrderID      *int64    `json:"order_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (CoinTransaction) TableName() string { return "coin_transactions" }

// VIP

type VipPackage struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string    `json:"name" gorm:"size:50;not null"`
	Price         int       `json:"price" gorm:"not null"`
	OriginalPrice int       `json:"original_price"`
	DurationDays  *int      `json:"duration_days"` // nil = permanent
	Benefits      string    `json:"benefits" gorm:"type:json"`
	IsActive      bool      `json:"is_active" gorm:"not null;default:true"`
	SortOrder     int       `json:"sort_order" gorm:"not null;default:0"`
	CreatedAt     time.Time `json:"created_at"`
}

func (VipPackage) TableName() string { return "vip_packages" }

type UserVip struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64      `json:"user_id" gorm:"not null"`
	PackageID int        `json:"package_id" gorm:"not null"`
	StartedAt time.Time  `json:"started_at"`
	ExpiresAt *time.Time `json:"expires_at"`
	IsActive  bool       `json:"is_active" gorm:"not null;default:true"`
	CreatedAt time.Time  `json:"created_at"`
}

func (UserVip) TableName() string { return "user_vip" }

// Order

type Order struct {
	ID         int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderNo    string     `json:"order_no" gorm:"size:32;not null;uniqueIndex"`
	UserID     int64      `json:"user_id" gorm:"not null"`
	Type       string     `json:"type" gorm:"size:20;not null"` // course, vip, coin
	TargetID   *int       `json:"target_id"`
	TargetName string     `json:"target_name" gorm:"size:200"`
	Amount     int        `json:"amount" gorm:"not null"`
	PayMethod  string     `json:"pay_method" gorm:"size:20"`
	Status     string     `json:"status" gorm:"size:20;not null;default:pending"`
	PaidAt     *time.Time `json:"paid_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (Order) TableName() string { return "orders" }

// Purchase

type Purchase struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id" gorm:"not null"`
	CourseID  int64     `json:"course_id" gorm:"not null"`
	OrderID   *int64    `json:"order_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (Purchase) TableName() string { return "purchases" }

// Guestbook

type Guestbook struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	LikeCount int       `json:"like_count" gorm:"not null;default:0"`
	IsPinned  bool      `json:"is_pinned" gorm:"not null;default:false"`
	Status    string    `json:"status" gorm:"size:20;not null;default:visible"`
	CreatedAt time.Time `json:"created_at"`

	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (Guestbook) TableName() string { return "guestbook" }

type GuestbookLike struct {
	UserID      int64 `gorm:"primaryKey"`
	GuestbookID int64 `gorm:"primaryKey"`
}

func (GuestbookLike) TableName() string { return "guestbook_likes" }

// Comment

type Comment struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     int64     `json:"user_id" gorm:"not null"`
	TargetType string    `json:"target_type" gorm:"size:20;not null"`
	TargetID   int64     `json:"target_id" gorm:"not null"`
	ParentID   *int64    `json:"parent_id"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	Status     string    `json:"status" gorm:"size:20;not null;default:visible"`
	CreatedAt  time.Time `json:"created_at"`

	User     *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Children []Comment `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

func (Comment) TableName() string { return "comments" }

// Home Config

type NavItem struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Title      string    `json:"title" gorm:"size:50;not null"`
	Subtitle   string    `json:"subtitle" gorm:"size:100"`
	Icon       string    `json:"icon" gorm:"size:200"`
	URL        string    `json:"url" gorm:"size:500;not null"`
	CategoryID *int      `json:"category_id"`
	SortOrder  int       `json:"sort_order" gorm:"not null;default:0"`
	IsActive   bool      `json:"is_active" gorm:"not null;default:true"`
	CreatedAt  time.Time `json:"created_at"`
}

func (NavItem) TableName() string { return "nav_items" }

type Banner struct {
	ID              int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Title           string    `json:"title" gorm:"size:100"`
	ImageURL        string    `json:"image_url" gorm:"size:500;not null"`
	LinkURL         string    `json:"link_url" gorm:"size:500"`
	Placement       string    `json:"placement" gorm:"size:50;not null;default:home"`
	BackgroundColor string    `json:"background_color" gorm:"size:30"`
	SortOrder       int       `json:"sort_order" gorm:"not null;default:0"`
	IsActive        bool      `json:"is_active" gorm:"not null;default:true"`
	CreatedAt       time.Time `json:"created_at"`
}

func (Banner) TableName() string { return "banners" }

// UserDownload

type UserDownload struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id" gorm:"not null"`
	CourseID  int64     `json:"course_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`

	Course *Course `json:"course,omitempty" gorm:"foreignKey:CourseID"`
}

func (UserDownload) TableName() string { return "user_downloads" }
