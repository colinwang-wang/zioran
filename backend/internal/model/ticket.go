package model

import "time"

// Ticket

type Ticket struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id" gorm:"not null"`
	Title     string    `json:"title" gorm:"size:200;not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	Status    string    `json:"status" gorm:"size:20;not null;default:open"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User    *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Replies []TicketReply  `json:"replies,omitempty" gorm:"foreignKey:TicketID"`
}

func (Ticket) TableName() string { return "tickets" }

type TicketReply struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TicketID  int64     `json:"ticket_id" gorm:"not null"`
	UserID    int64     `json:"user_id" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	IsAdmin   bool      `json:"is_admin" gorm:"not null;default:false"`
	CreatedAt time.Time `json:"created_at"`

	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (TicketReply) TableName() string { return "ticket_replies" }

// Setting

type Setting struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Key       string    `json:"key" gorm:"size:100;not null;uniqueIndex"`
	Value     string    `json:"value" gorm:"type:text"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Setting) TableName() string { return "settings" }

// OperationLog

type OperationLog struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	AdminID   int64     `json:"admin_id" gorm:"not null"`
	Action    string    `json:"action" gorm:"size:100;not null"`
	Target    string    `json:"target" gorm:"size:200"`
	Detail    string    `json:"detail" gorm:"type:text"`
	IP        string    `json:"ip" gorm:"size:50"`
	CreatedAt time.Time `json:"created_at"`
}

func (OperationLog) TableName() string { return "operation_logs" }

// PaymentLog

type PaymentLog struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID   *int64    `json:"order_id"`
	Type      string    `json:"type" gorm:"size:50;not null"`
	Detail    string    `json:"detail" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
}

func (PaymentLog) TableName() string { return "payment_logs" }

// DTOs

type CreateTicketRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type TicketReplyRequest struct {
	Content string `json:"content" binding:"required"`
}

type TicketStatusRequest struct {
	Status string `json:"status" binding:"required"` // processing, replied, closed
}

type TicketResponse struct {
	ID        int64               `json:"id"`
	UserID    int64               `json:"user_id"`
	Username  string              `json:"username"`
	Title     string              `json:"title"`
	Content   string              `json:"content"`
	Status    string              `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
	Replies   []TicketReplyResponse `json:"replies,omitempty"`
}

type TicketReplyResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

type ForgotPasswordRequest struct {
	Phone       string `json:"phone" binding:"required"`
	SMSCode     string `json:"sms_code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type AdminCreateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required"`
}

type AdminUpdateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type SettingsMap map[string]string

type FinanceSummary struct {
	TodayRevenue   int64 `json:"today_revenue"`
	TotalSettled   int64 `json:"total_settled"`
	TotalPending   int64 `json:"total_pending"`
}

type CommentReplyRequest struct {
	Content string `json:"content" binding:"required"`
}
