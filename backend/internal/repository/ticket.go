package repository

import (
	"context"
	"time"

	"github.com/zioran/backend/internal/model"
	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

// Tickets

func (r *TicketRepository) Create(ctx context.Context, ticket *model.Ticket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

func (r *TicketRepository) FindByID(ctx context.Context, id int64) (*model.Ticket, error) {
	var t model.Ticket
	err := r.db.WithContext(ctx).Preload("Replies").Preload("Replies.User").Preload("User").Preload("Attachments").First(&t, id).Error
	return &t, err
}

func (r *TicketRepository) ListByUser(ctx context.Context, userID int64, page, pageSize int) ([]model.Ticket, int64, error) {
	var items []model.Ticket
	var total int64
	q := r.db.WithContext(ctx).Model(&model.Ticket{}).Where("user_id = ?", userID)
	q.Count(&total)
	err := q.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *TicketRepository) ListAll(ctx context.Context, page, pageSize int, status string) ([]model.Ticket, int64, error) {
	var items []model.Ticket
	var total int64
	q := r.db.WithContext(ctx).Model(&model.Ticket{}).Preload("User")
	if status != "" {
		q = q.Where("status = ?", normalizeTicketStatus(status))
	}
	q.Count(&total)
	err := q.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *TicketRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&model.Ticket{}).Where("id = ?", id).Update("status", status).Error
}

func (r *TicketRepository) CreateReply(ctx context.Context, reply *model.TicketReply) error {
	return r.db.WithContext(ctx).Create(reply).Error
}

func (r *TicketRepository) CreateAttachments(ctx context.Context, attachments []model.TicketAttachment) error {
	if len(attachments) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&attachments).Error
}

// Settings

func (r *TicketRepository) GetAllSettings(ctx context.Context) ([]model.Setting, error) {
	var settings []model.Setting
	err := r.db.WithContext(ctx).Find(&settings).Error
	return settings, err
}

func (r *TicketRepository) UpsertSettings(ctx context.Context, settings model.SettingsMap) error {
	for k, v := range settings {
		r.db.WithContext(ctx).Where(model.Setting{Key: k}).Assign(model.Setting{Value: v}).FirstOrCreate(&model.Setting{})
	}
	return nil
}

// Admins

func (r *TicketRepository) ListAdmins(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := r.db.WithContext(ctx).Where("role = ?", "admin").Find(&users).Error
	return users, err
}

func (r *TicketRepository) CreateAdmin(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *TicketRepository) UpdateAdmin(ctx context.Context, id int64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *TicketRepository) DeleteAdmin(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// Logs

func (r *TicketRepository) OperationLogs(ctx context.Context, page, pageSize int) ([]model.OperationLog, int64, error) {
	var items []model.OperationLog
	var total int64
	q := r.db.WithContext(ctx).Model(&model.OperationLog{})
	q.Count(&total)
	err := q.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *TicketRepository) PaymentLogs(ctx context.Context, page, pageSize int) ([]model.PaymentLog, int64, error) {
	var items []model.PaymentLog
	var total int64
	q := r.db.WithContext(ctx).Model(&model.PaymentLog{})
	q.Count(&total)
	err := q.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

// Finance

func (r *TicketRepository) FinanceSummary(ctx context.Context) (*model.FinanceSummary, error) {
	var todayRevenue, totalSettled, totalPending int64
	todayStart := startOfToday()
	tomorrowStart := todayStart.AddDate(0, 0, 1)
	r.db.WithContext(ctx).Model(&model.Order{}).Where("status = ? AND paid_at >= ? AND paid_at < ?", "paid", todayStart, tomorrowStart).Select("COALESCE(SUM(amount),0)").Scan(&todayRevenue)
	r.db.WithContext(ctx).Model(&model.Order{}).Where("status = ?", "paid").Select("COALESCE(SUM(amount),0)").Scan(&totalSettled)
	r.db.WithContext(ctx).Model(&model.Order{}).Where("status = ?", "pending").Select("COALESCE(SUM(amount),0)").Scan(&totalPending)
	return &model.FinanceSummary{TodayRevenue: todayRevenue, TotalSettled: totalSettled, TotalPending: totalPending}, nil
}

func (r *TicketRepository) FinanceWithdrawals(ctx context.Context, page, pageSize int, status string) ([]model.WithdrawalRequest, int64, error) {
	var items []model.WithdrawalRequest
	var total int64
	q := r.db.WithContext(ctx).Model(&model.WithdrawalRequest{}).Preload("User")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	q.Count(&total)
	err := q.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

// Comments (for admin reply)

func (r *TicketRepository) CreateComment(ctx context.Context, comment *model.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *TicketRepository) FindComment(ctx context.Context, id int64) (*model.Comment, error) {
	var c model.Comment
	err := r.db.WithContext(ctx).First(&c, id).Error
	return &c, err
}

// Orders

func (r *TicketRepository) GetOrder(ctx context.Context, id int64) (*model.Order, error) {
	var o model.Order
	err := r.db.WithContext(ctx).First(&o, id).Error
	return &o, err
}

func (r *TicketRepository) UpdateOrderStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error
}

func normalizeTicketStatus(status string) string {
	if status == "pending" {
		return "open"
	}
	return status
}

func startOfToday() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}
