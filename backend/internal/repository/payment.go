package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zioran/backend/internal/model"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// Coin Account

func (r *PaymentRepository) GetOrCreateAccount(ctx context.Context, userID int64) (*model.CoinAccount, error) {
	var acc model.CoinAccount
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&acc).Error
	if err == gorm.ErrRecordNotFound {
		acc = model.CoinAccount{UserID: userID}
		err = r.db.WithContext(ctx).Create(&acc).Error
	}
	return &acc, err
}

func (r *PaymentRepository) CoinTransactions(ctx context.Context, userID int64, page, pageSize int) ([]model.CoinTransaction, int64, error) {
	var txs []model.CoinTransaction
	var total int64
	query := r.db.WithContext(ctx).Model(&model.CoinTransaction{}).Where("user_id = ?", userID)
	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&txs).Error
	return txs, total, err
}

// Recharge (with transaction safety)
func (r *PaymentRepository) Recharge(ctx context.Context, userID int64, amount int, orderID int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var acc model.CoinAccount
		if err := tx.Where("user_id = ?", userID).First(&acc).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				acc = model.CoinAccount{UserID: userID}
				tx.Create(&acc)
			} else {
				return err
			}
		}
		acc.Balance += amount
		acc.TotalEarned += amount
		if err := tx.Save(&acc).Error; err != nil {
			return err
		}
		txn := model.CoinTransaction{
			UserID:       userID,
			Type:         "recharge",
			Amount:       amount,
			BalanceAfter: acc.Balance,
			Description:  fmt.Sprintf("充值%d金币", amount),
			OrderID:      &orderID,
		}
		return tx.Create(&txn).Error
	})
}

// Deduct coins (with balance check in transaction)
func (r *PaymentRepository) DeductCoins(ctx context.Context, userID int64, amount int, txType, desc string, orderID *int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var acc model.CoinAccount
		// Lock row for update
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("user_id = ?", userID).First(&acc).Error; err != nil {
			return fmt.Errorf("account not found")
		}
		if acc.Balance < amount {
			return fmt.Errorf("insufficient balance")
		}
		acc.Balance -= amount
		acc.TotalSpent += amount
		if err := tx.Save(&acc).Error; err != nil {
			return err
		}
		txn := model.CoinTransaction{
			UserID:       userID,
			Type:         txType,
			Amount:       -amount,
			BalanceAfter: acc.Balance,
			Description:  desc,
			OrderID:      orderID,
		}
		return tx.Create(&txn).Error
	})
}

// VIP

func (r *PaymentRepository) VipPackages(ctx context.Context) ([]model.VipPackage, error) {
	var pkgs []model.VipPackage
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("sort_order ASC").Find(&pkgs).Error
	return pkgs, err
}

func (r *PaymentRepository) GetVipStatus(ctx context.Context, userID int64) (*model.UserVip, error) {
	var vip model.UserVip
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true).
		Where("expires_at IS NULL OR expires_at > ?", time.Now()).
		Order("created_at DESC").First(&vip).Error
	if err != nil {
		return nil, err
	}
	return &vip, nil
}

func (r *PaymentRepository) CreateVip(ctx context.Context, vip *model.UserVip) error {
	return r.db.WithContext(ctx).Create(vip).Error
}

func (r *PaymentRepository) GetVipPackage(ctx context.Context, id int) (*model.VipPackage, error) {
	var pkg model.VipPackage
	err := r.db.WithContext(ctx).First(&pkg, id).Error
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

// Orders

func (r *PaymentRepository) CreateOrder(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *PaymentRepository) GetOrder(ctx context.Context, id int64) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *PaymentRepository) GetOrderByNo(ctx context.Context, orderNo string) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *PaymentRepository) UpdateOrderStatus(ctx context.Context, id int64, status string) error {
	updates := map[string]interface{}{"status": status}
	if status == "paid" {
		now := time.Now()
		updates["paid_at"] = &now
	}
	return r.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PaymentRepository) UserOrders(ctx context.Context, userID int64, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", userID)
	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders).Error
	return orders, total, err
}

func (r *PaymentRepository) AdminOrders(ctx context.Context, page, pageSize int, status string) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Order{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders).Error
	return orders, total, err
}

// Purchases

func (r *PaymentRepository) HasPurchased(ctx context.Context, userID, courseID int64) bool {
	var count int64
	r.db.WithContext(ctx).Model(&model.Purchase{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).Count(&count)
	return count > 0
}

func (r *PaymentRepository) CreatePurchase(ctx context.Context, p *model.Purchase) error {
	return r.db.WithContext(ctx).Create(p).Error
}

// Downloads

func (r *PaymentRepository) RecordDownload(ctx context.Context, userID, courseID int64) error {
	return r.db.WithContext(ctx).Create(&model.UserDownload{UserID: userID, CourseID: courseID}).Error
}

func (r *PaymentRepository) UserDownloads(ctx context.Context, userID int64, page, pageSize int) ([]model.UserDownload, int64, error) {
	var downloads []model.UserDownload
	var total int64
	query := r.db.WithContext(ctx).Model(&model.UserDownload{}).Where("user_id = ?", userID)
	query.Count(&total)
	err := query.Preload("Course").Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&downloads).Error
	return downloads, total, err
}

// User favorites (paginated)

func (r *PaymentRepository) UserFavorites(ctx context.Context, userID int64, page, pageSize int) ([]model.UserFavorite, int64, error) {
	var favs []model.UserFavorite
	var total int64
	query := r.db.WithContext(ctx).Model(&model.UserFavorite{}).Where("user_id = ?", userID)
	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&favs).Error
	return favs, total, err
}

func (r *PaymentRepository) GetCoursesByIDs(ctx context.Context, ids []int64) ([]model.Course, error) {
	var courses []model.Course
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&courses).Error
	return courses, err
}

// IsVip check helper
func (r *PaymentRepository) IsVip(ctx context.Context, userID int64) bool {
	_, err := r.GetVipStatus(ctx, userID)
	return err == nil
}

// Favorites

func (r *PaymentRepository) AddFavorite(ctx context.Context, userID, courseID int64) error {
	fav := model.UserFavorite{UserID: userID, CourseID: courseID}
	return r.db.WithContext(ctx).FirstOrCreate(&fav, fav).Error
}

func (r *PaymentRepository) RemoveFavorite(ctx context.Context, userID, courseID int64) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND course_id = ?", userID, courseID).Delete(&model.UserFavorite{}).Error
}

// Admin users
func (r *PaymentRepository) AdminUsers(ctx context.Context, page, pageSize int, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	query := r.db.WithContext(ctx).Model(&model.User{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	return users, total, err
}

func (r *PaymentRepository) UpdateUserStatus(ctx context.Context, userID int64, status string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("status", status).Error
}

// Dashboard
func (r *PaymentRepository) DashboardStats(ctx context.Context) (*model.DashboardStats, error) {
	var stats model.DashboardStats
	r.db.WithContext(ctx).Model(&model.User{}).Count(&stats.TotalUsers)
	r.db.WithContext(ctx).Model(&model.Course{}).Count(&stats.TotalCourses)
	r.db.WithContext(ctx).Model(&model.Order{}).Count(&stats.TotalOrders)

	var todayRevenue *int64
	today := time.Now().Format("2006-01-02")
	r.db.WithContext(ctx).Model(&model.Order{}).Where("status = ? AND DATE(created_at) = ?", "paid", today).
		Select("COALESCE(SUM(amount), 0)").Scan(&todayRevenue)
	if todayRevenue != nil {
		stats.TodayRevenue = *todayRevenue
	}

	// Growth rates (compare today vs yesterday counts)
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	var todayUsers, yesterdayUsers int64
	r.db.WithContext(ctx).Model(&model.User{}).Where("DATE(created_at) = ?", today).Count(&todayUsers)
	r.db.WithContext(ctx).Model(&model.User{}).Where("DATE(created_at) = ?", yesterday).Count(&yesterdayUsers)
	stats.UserGrowth = calcGrowth(todayUsers, yesterdayUsers)

	var todayOrders, yesterdayOrders int64
	r.db.WithContext(ctx).Model(&model.Order{}).Where("DATE(created_at) = ?", today).Count(&todayOrders)
	r.db.WithContext(ctx).Model(&model.Order{}).Where("DATE(created_at) = ?", yesterday).Count(&yesterdayOrders)
	stats.OrderGrowth = calcGrowth(todayOrders, yesterdayOrders)

	var todayCourses, yesterdayCourses int64
	r.db.WithContext(ctx).Model(&model.Course{}).Where("DATE(created_at) = ?", today).Count(&todayCourses)
	r.db.WithContext(ctx).Model(&model.Course{}).Where("DATE(created_at) = ?", yesterday).Count(&yesterdayCourses)
	stats.CourseGrowth = calcGrowth(todayCourses, yesterdayCourses)

	var yesterdayRevenue int64
	r.db.WithContext(ctx).Model(&model.Order{}).Where("status = ? AND DATE(created_at) = ?", "paid", yesterday).
		Select("COALESCE(SUM(amount), 0)").Scan(&yesterdayRevenue)
	stats.RevenueGrowth = calcGrowth(stats.TodayRevenue, yesterdayRevenue)

	return &stats, nil
}

func calcGrowth(current, previous int64) float64 {
	if previous == 0 {
		if current > 0 {
			return 100
		}
		return 0
	}
	return float64(current-previous) / float64(previous) * 100
}
