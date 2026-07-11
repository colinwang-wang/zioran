package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zioran/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

type UserAdminStats struct {
	Balance        int
	PurchasedCount int64
	FavoriteCount  int64
}

type DownloadOrderMeta struct {
	CourseID int64  `gorm:"column:course_id"`
	OrderNo  string `gorm:"column:order_no"`
	Amount   int    `gorm:"column:amount"`
}

func (r *PaymentRepository) GetSetting(ctx context.Context, key string) (string, error) {
	var setting model.Setting
	err := r.db.WithContext(ctx).Where("`key` = ?", key).First(&setting).Error
	if err != nil {
		return "", err
	}
	return setting.Value, nil
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

func (r *PaymentRepository) AdminVipPackages(ctx context.Context) ([]model.VipPackage, error) {
	var pkgs []model.VipPackage
	err := r.db.WithContext(ctx).Order("sort_order ASC, id ASC").Find(&pkgs).Error
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

func (r *PaymentRepository) UpdateVipPackage(ctx context.Context, id int, updates map[string]interface{}) (*model.VipPackage, error) {
	if err := r.db.WithContext(ctx).Model(&model.VipPackage{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.GetVipPackage(ctx, id)
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

func (r *PaymentRepository) CompleteRechargeOrder(ctx context.Context, orderID int64, coins int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order model.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, orderID).Error; err != nil {
			return err
		}
		if order.Status != "paid" {
			now := time.Now()
			if err := tx.Model(&model.Order{}).Where("id = ?", order.ID).
				Updates(map[string]interface{}{"status": "paid", "paid_at": &now}).Error; err != nil {
				return err
			}
		}

		var existing int64
		if err := tx.Model(&model.CoinTransaction{}).
			Where("order_id = ? AND type = ?", order.ID, "recharge").
			Count(&existing).Error; err != nil {
			return err
		}
		if existing > 0 {
			return nil
		}

		var acc model.CoinAccount
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", order.UserID).First(&acc).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				acc = model.CoinAccount{UserID: order.UserID}
				if createErr := tx.Create(&acc).Error; createErr != nil {
					return createErr
				}
			} else {
				return err
			}
		}
		acc.Balance += coins
		acc.TotalEarned += coins
		if err := tx.Save(&acc).Error; err != nil {
			return err
		}
		txn := model.CoinTransaction{
			UserID:       order.UserID,
			Type:         "recharge",
			Amount:       coins,
			BalanceAfter: acc.Balance,
			Description:  fmt.Sprintf("充值%d金币", coins),
			OrderID:      &order.ID,
		}
		return tx.Create(&txn).Error
	})
}

func (r *PaymentRepository) UpdateOrderStatus(ctx context.Context, id int64, status string) error {
	updates := map[string]interface{}{"status": status}
	if status == "paid" {
		now := time.Now()
		updates["paid_at"] = &now
	}
	return r.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PaymentRepository) RefundOrder(ctx context.Context, orderID int64, coinRefundLimit int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order model.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, orderID).Error; err != nil {
			return err
		}
		if order.Status != "paid" {
			return fmt.Errorf("order status does not allow refund")
		}

		switch order.Type {
		case "coin":
			var acc model.CoinAccount
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", order.UserID).First(&acc).Error; err != nil {
				return err
			}
			refundCoins := coinRefundLimit
			if acc.Balance < refundCoins {
				refundCoins = acc.Balance
			}
			if refundCoins <= 0 {
				return fmt.Errorf("insufficient balance for refund")
			}
			acc.Balance -= refundCoins
			acc.TotalSpent += refundCoins
			if err := tx.Save(&acc).Error; err != nil {
				return err
			}
			txn := model.CoinTransaction{
				UserID:       order.UserID,
				Type:         "refund",
				Amount:       -refundCoins,
				BalanceAfter: acc.Balance,
				Description:  "退款: " + order.TargetName,
				OrderID:      &order.ID,
			}
			if err := tx.Create(&txn).Error; err != nil {
				return err
			}
		case "course", "vip":
			var acc model.CoinAccount
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", order.UserID).First(&acc).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					acc = model.CoinAccount{UserID: order.UserID}
					if createErr := tx.Create(&acc).Error; createErr != nil {
						return createErr
					}
				} else {
					return err
				}
			}
			acc.Balance += order.Amount
			acc.TotalEarned += order.Amount
			if err := tx.Save(&acc).Error; err != nil {
				return err
			}
			txn := model.CoinTransaction{
				UserID:       order.UserID,
				Type:         "refund",
				Amount:       order.Amount,
				BalanceAfter: acc.Balance,
				Description:  "退款: " + order.TargetName,
				OrderID:      &order.ID,
			}
			if err := tx.Create(&txn).Error; err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported order type")
		}

		// 退款后撤销下载权限：删除该订单对应的 purchase 记录
		if order.Type == "course" {
			if err := tx.Where("order_id = ?", order.ID).Delete(&model.Purchase{}).Error; err != nil {
				return err
			}
		}

		return tx.Model(&model.Order{}).Where("id = ? AND status = ?", order.ID, "paid").Update("status", "refunded").Error
	})
}

func (r *PaymentRepository) LogPayment(ctx context.Context, orderID *int64, logType, detail string) error {
	log := &model.PaymentLog{OrderID: orderID, Type: logType, Detail: detail}
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *PaymentRepository) UserOrders(ctx context.Context, userID int64, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", userID)
	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders).Error
	return orders, total, err
}

func (r *PaymentRepository) AdminOrders(ctx context.Context, page, pageSize int, filter model.AdminOrderFilter) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Order{})
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", normalizeOrderType(filter.Type))
	}
	if filter.Keyword != "" {
		kw := "%" + filter.Keyword + "%"
		query = query.Where(
			"(order_no LIKE ? OR target_name LIKE ? OR user_id IN (SELECT id FROM users WHERE username LIKE ? OR phone LIKE ?))",
			kw, kw, kw, kw,
		)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at < ?", *filter.EndDate)
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

func (r *PaymentRepository) UserDownloadOrderMetas(ctx context.Context, userID int64, courseIDs []int64) (map[int64]DownloadOrderMeta, error) {
	metas := make(map[int64]DownloadOrderMeta, len(courseIDs))
	if len(courseIDs) == 0 {
		return metas, nil
	}
	var rows []DownloadOrderMeta
	err := r.db.WithContext(ctx).
		Table("purchases").
		Select("purchases.course_id, orders.order_no, orders.amount").
		Joins("LEFT JOIN orders ON orders.id = purchases.order_id").
		Where("purchases.user_id = ? AND purchases.course_id IN ?", userID, courseIDs).
		Order("purchases.created_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		if _, exists := metas[row.CourseID]; !exists {
			metas[row.CourseID] = row
		}
	}
	return metas, nil
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
func (r *PaymentRepository) AdminUsers(ctx context.Context, page, pageSize int, keyword string, vipFilter string) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	query := r.db.WithContext(ctx).Model(&model.User{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	switch vipFilter {
	case "vip":
		query = query.Where("id IN (SELECT user_id FROM user_vip WHERE is_active = ? AND (expires_at IS NULL OR expires_at > ?))", true, time.Now())
	case "normal":
		query = query.Where("id NOT IN (SELECT user_id FROM user_vip WHERE is_active = ? AND (expires_at IS NULL OR expires_at > ?))", true, time.Now())
	}
	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	return users, total, err
}

func (r *PaymentRepository) UserAdminStats(ctx context.Context, userID int64) (UserAdminStats, error) {
	stats := UserAdminStats{}
	var acc model.CoinAccount
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&acc).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return stats, err
	}
	if err == nil {
		stats.Balance = acc.Balance
	}
	if err := r.db.WithContext(ctx).Model(&model.Purchase{}).
		Where("user_id = ?", userID).
		Distinct("course_id").
		Count(&stats.PurchasedCount).Error; err != nil {
		return stats, err
	}
	if err := r.db.WithContext(ctx).Model(&model.UserFavorite{}).
		Where("user_id = ?", userID).
		Count(&stats.FavoriteCount).Error; err != nil {
		return stats, err
	}
	return stats, nil
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

	// 新增统计
	r.db.WithContext(ctx).Model(&model.UserFavorite{}).Count(&stats.TotalFavorites)
	r.db.WithContext(ctx).Model(&model.Order{}).Where("status = ?", "pending").Count(&stats.PendingOrders)
	threeDaysAgo := time.Now().AddDate(0, 0, -3)
	r.db.WithContext(ctx).Model(&model.User{}).Where("created_at >= ?", threeDaysAgo).Count(&stats.RecentNewUsers)

	var todayRevenue *int64
	now := time.Now()
	todayStart := startOfDay(now)
	tomorrowStart := todayStart.AddDate(0, 0, 1)
	yesterdayStart := todayStart.AddDate(0, 0, -1)
	r.db.WithContext(ctx).Model(&model.Order{}).Where("status = ? AND created_at >= ? AND created_at < ?", "paid", todayStart, tomorrowStart).
		Select("COALESCE(SUM(amount), 0)").Scan(&todayRevenue)
	if todayRevenue != nil {
		stats.TodayRevenue = *todayRevenue
	}

	// Growth rates (compare today vs yesterday counts)
	var todayUsers, yesterdayUsers int64
	r.db.WithContext(ctx).Model(&model.User{}).Where("created_at >= ? AND created_at < ?", todayStart, tomorrowStart).Count(&todayUsers)
	r.db.WithContext(ctx).Model(&model.User{}).Where("created_at >= ? AND created_at < ?", yesterdayStart, todayStart).Count(&yesterdayUsers)
	stats.UserGrowth = calcGrowth(todayUsers, yesterdayUsers)

	var todayOrders, yesterdayOrders int64
	r.db.WithContext(ctx).Model(&model.Order{}).Where("created_at >= ? AND created_at < ?", todayStart, tomorrowStart).Count(&todayOrders)
	r.db.WithContext(ctx).Model(&model.Order{}).Where("created_at >= ? AND created_at < ?", yesterdayStart, todayStart).Count(&yesterdayOrders)
	stats.OrderGrowth = calcGrowth(todayOrders, yesterdayOrders)

	var todayCourses, yesterdayCourses int64
	r.db.WithContext(ctx).Model(&model.Course{}).Where("created_at >= ? AND created_at < ?", todayStart, tomorrowStart).Count(&todayCourses)
	r.db.WithContext(ctx).Model(&model.Course{}).Where("created_at >= ? AND created_at < ?", yesterdayStart, todayStart).Count(&yesterdayCourses)
	stats.CourseGrowth = calcGrowth(todayCourses, yesterdayCourses)

	var yesterdayRevenue int64
	r.db.WithContext(ctx).Model(&model.Order{}).Where("status = ? AND created_at >= ? AND created_at < ?", "paid", yesterdayStart, todayStart).
		Select("COALESCE(SUM(amount), 0)").Scan(&yesterdayRevenue)
	stats.RevenueGrowth = calcGrowth(stats.TodayRevenue, yesterdayRevenue)

	return &stats, nil
}

func (r *PaymentRepository) DashboardDailyCounts(ctx context.Context, start, end time.Time) (int64, int64, error) {
	var users, orders int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("created_at >= ? AND created_at < ?", start, end).Count(&users).Error; err != nil {
		return 0, 0, err
	}
	if err := r.db.WithContext(ctx).Model(&model.Order{}).Where("created_at >= ? AND created_at < ?", start, end).Count(&orders).Error; err != nil {
		return 0, 0, err
	}
	return users, orders, nil
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

func normalizeOrderType(orderType string) string {
	switch orderType {
	case "coin_recharge":
		return "coin"
	case "vip_purchase":
		return "vip"
	case "course_purchase":
		return "course"
	default:
		return orderType
	}
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
